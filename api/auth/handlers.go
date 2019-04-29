package auth

import (
	"net/http"
	"regexp"

	"github.com/dbeliakov/revisor/api/store"
	"github.com/dbeliakov/revisor/api/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

// APIUser represents api result struct
type APIUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Login     string `json:"username"`
}

// NewAPIUser creates new api user from store user
func NewAPIUser(u store.User) APIUser {
	return APIUser{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Login:     u.Login,
	}
}

var checkLogin = regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`).MatchString

// LoginHandler checks username and password and returns jwt on success
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var form struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	if err := utils.UnmarshalForm(w, r, &form); err != nil {
		return
	}

	user, err := store.Auth.FindUserByLogin(form.Username)
	if err != nil {
		logrus.Infof("Cannot find user with such login: %s, error: %+v", form.Username, err)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusNotAcceptable,
			Message:       "Incorrect username field",
			ClientMessage: "Пользователь с таким логином не найден",
		})
		return
	}
	if checkPassword(user, form.Password) {
		token, err := newToken(user)
		if err != nil {
			logrus.Errorf("Cannot create new token for user: %s, error: %+v", form.Username, err)
			utils.Error(w, utils.InternalErrorResponse("Cannot generate token"))
			return
		}
		w.Header().Add("Authorization", "Bearer "+token)
		utils.Ok(w, nil)
	} else {
		logrus.Infof("Invalid password for user: %s", form.Username)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusNotAcceptable,
			Message:       "Incorrect password",
			ClientMessage: "Неправильный пароль",
		})
	}
}

// SignUpHandler creates new user
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var form struct {
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
		Username  string `json:"username" validate:"required"`
		Password  string `json:"password" validate:"required,min=6"`
	}
	if err := utils.UnmarshalForm(w, r, &form); err != nil {
		return
	}

	if !checkLogin(form.Username) {
		logrus.Infof("Incorrect login: %s", form.Username)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusConflict,
			Message:       "Incorrect login",
			ClientMessage: "В логине могут быть только латинские символы, цифры и знаки '-' и '_'",
		})
		return
	}

	user, err := newUser(form.FirstName, form.LastName, form.Username, form.Password)
	if err != nil {
		logrus.Warnf("Cannot create new user: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot create new user"))
		return
	}
	err = store.Auth.CreateUser(user)
	if xerrors.Is(err, store.ErrUserExists) {
		logrus.Infof("Login is not free: %s", form.Username)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusConflict,
			Message:       "Username is not free",
			ClientMessage: "Пользователь с таким логином уже зарегистрирован",
		})
		return
	} else if err != nil {
		logrus.Warnf("Cannot save new user to database: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("Cannot create new user"))
		return
	}

	utils.Ok(w, nil)
}

// UserInfoHandler returns info about user
var UserInfoHandler = AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	u, err := UserFromRequest(r)
	if err != nil {
		logrus.Errorf("Error while getting user from request context: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("No authorized user for this request"))
		return
	}
	utils.Ok(w, NewAPIUser(u))
})

// ChangePasswordHandler changes user password
var ChangePasswordHandler = AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	user, err := UserFromRequest(r)
	if err != nil {
		logrus.Errorf("Error while getting user from request context: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("No authorized user for this request"))
		return
	}

	var form struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=6"`
	}
	if err := utils.UnmarshalForm(w, r, &form); err != nil {
		return
	}

	if checkPassword(user, form.OldPassword) {
		err := setPassword(&user, form.NewPassword)
		if err != nil {
			logrus.Errorf("Cannot set new password for user: %s, error: %+v", user.Login, err)
			utils.Error(w, utils.InternalErrorResponse("Cannot store password"))
			return
		}
		err = store.Auth.UpdateUser(user)
		if err != nil {
			logrus.Errorf("Cannot save new password for user: %s, error: %+v", user.Login, err)
			utils.Error(w, utils.InternalErrorResponse("Cannot store password"))
			return
		}
		utils.Ok(w, nil)
	} else {
		logrus.Infof("Cannot change password for user %s: incorrect current password", user.Login)
		utils.Error(w, utils.JSONErrorResponse{
			Status:        http.StatusNotAcceptable,
			Message:       "Incorrect current password",
			ClientMessage: "Неправильный текущий пароль",
		})
	}
})
