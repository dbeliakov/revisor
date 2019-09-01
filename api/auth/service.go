package auth

import (
	"net/http"
	"regexp"
	"time"

	"go.uber.org/zap"
	"golang.org/x/xerrors"

	"github.com/dbeliakov/revisor/api/store"
	"github.com/dbeliakov/revisor/api/utils"
)

const (
	cookieDuration = 31 * 24 * time.Hour
	cookieMaxAge   = time.Duration(365 * 24 * time.Hour)
)

// Service provides http handlers for auth
type Service struct {
	logger *zap.Logger
	store  store.AuthStore
}

// NewService creates new auth service
func NewService(logger *zap.Logger, store store.AuthStore) *Service {
	return &Service{
		logger: logger,
		store:  store,
	}
}

// HandleLogin is http handler that creates JWT token by correct username and password
func (s *Service) HandleLogin(rw http.ResponseWriter, r *http.Request) {
	var form struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	if err := utils.JSONUnmarshalAndValidate(r.Body, &form); err != nil {
		s.logger.Warn("Cannot unmarshal request body", zap.Error(err))
		utils.HTTPBadRequest(rw)
		return
	}

	user, err := s.store.FindUserByLogin(form.Username)
	if err != nil {
		s.logger.Warn("Cannot find user", zap.String("username", form.Username), zap.Error(err))
		utils.HTTPError(rw, http.StatusNotAcceptable, "Пользователь с таким логином не найден")
		return
	}

	if !checkPassword(user, form.Password) {
		s.logger.Info("Invalid password", zap.String("username", form.Username))
		utils.HTTPError(rw, http.StatusNotAcceptable, "Неправильный пароль")
		return
	}

	token, err := newToken(user)
	if err != nil {
		s.logger.Error("Cannot create new token for user", zap.String("username", form.Username), zap.Error(err))
		utils.HTTPInternalServerError(rw)
		return
	}
	http.SetCookie(rw, &http.Cookie{
		Name:     "revisor-token",
		Value:    token,
		Expires:  time.Now().Add(cookieDuration),
		Path:     "/",
		MaxAge:   int(cookieMaxAge.Seconds()),
		Secure:   false,
		HttpOnly: true,
	})
	utils.HTTPOK(rw, nil)
}

// HandleSignUp is http handler that creates new user and returns JWT token
func (s *Service) HandleSignUp(rw http.ResponseWriter, r *http.Request) {
	var form struct {
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
		Username  string `json:"username" validate:"required,min=3"`
		Password  string `json:"password" validate:"required,min=6"`
	}
	if err := utils.JSONUnmarshalAndValidate(r.Body, &form); err != nil {
		s.logger.Warn("Cannot unmarshal request body", zap.Error(err))
		utils.HTTPBadRequest(rw)
		return
	}

	var checkLogin = regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`).MatchString
	if !checkLogin(form.Username) {
		s.logger.Info("Incorrect login", zap.String("username", form.Username))
		utils.HTTPError(rw, http.StatusConflict,
			"В логине могут быть только латинские символы, цифры и знаки '-' и '_'")
		return
	}

	user, err := newUser(form.FirstName, form.LastName, form.Username, form.Password)
	if err != nil {
		s.logger.Error("Cannot create new user", zap.Error(err))
		utils.HTTPInternalServerError(rw)
		return
	}
	err = store.Auth.CreateUser(user)
	if xerrors.Is(err, store.ErrUserExists) {
		s.logger.Info("User already exists", zap.String("username", form.Username))
		utils.HTTPError(rw, http.StatusConflict, "Пользователь с таким логином уже зарегистрирован")
		return
	} else if err != nil {
		s.logger.Error("Cannot save new user to database", zap.Error(err))
		utils.HTTPInternalServerError(rw)
		return
	}

	utils.HTTPOK(rw, nil)
}

// HandleUserInfo is http handler that returns user info from JWT
// auth.Required middleware is required
func (s *Service) HandleUserInfo(rw http.ResponseWriter, r *http.Request) {
	u, err := UserFromRequest(r)
	if err != nil {
		s.logger.Error("Error while getting user from request context", zap.Error(err))
		utils.HTTPInternalServerError(rw)
		return
	}
	utils.Ok(rw, u)
}

// HandleChangePassword is http handler that chnages user password in database
func (s *Service) HandleChangePassword(rw http.ResponseWriter, r *http.Request) {
	user, err := UserFromRequest(r)
	if err != nil {
		s.logger.Error("Error while getting user from request context", zap.Error(err))
		utils.HTTPInternalServerError(rw)
		return
	}

	var form struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=6"`
	}
	if err := utils.JSONUnmarshalAndValidate(r.Body, &form); err != nil {
		s.logger.Warn("Cannot unmarshal request body", zap.Error(err))
		utils.HTTPBadRequest(rw)
		return
	}

	if !checkPassword(user, form.OldPassword) {
		s.logger.Info("Cannot change password for user: incorrect current password", zap.String("usearname", user.Login))
		utils.HTTPError(rw, http.StatusNotAcceptable, "Неправильный текущий пароль")
		return
	}

	err = setPassword(&user, form.NewPassword)
	if err != nil {
		s.logger.Error("Cannot set new password for user", zap.Error(err), zap.String("usearname", user.Login))
		utils.HTTPInternalServerError(rw)
		return
	}
	err = store.Auth.UpdateUser(user)
	if err != nil {
		s.logger.Error("Cannot save new password for user", zap.Error(err), zap.String("usearname", user.Login))
		utils.HTTPInternalServerError(rw)
		return
	}
	utils.Ok(rw, nil)
}
