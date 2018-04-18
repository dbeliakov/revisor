package handlers

import (
	"errors"
	"log"
	"net/http"
	"reviewer/api/database"
	"reviewer/api/middlewares"
	"reviewer/api/utils"
)

type loginForm struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginHandler checks username and password and returns jwt on success
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var form loginForm
	err := utils.UnmarshalForm(r, &form)
	if err == utils.ErrIncorrectBody {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	} else if err == utils.ErrIncorectFormFields {
		utils.Error(w, http.StatusNotAcceptable, err.Error())
		return
	} else if err != nil {
		log.Panic("Assert: unexpected error: " + err.Error())
	}

	user, err := database.UserByLogin(form.Username)
	if err != nil {
		utils.Error(w, http.StatusNotAcceptable, "Incorrect username field")
		return
	}
	if user.CheckPassword(form.Password) {
		token, err := user.NewToken(user.ID.Hex())
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Server error")
			return
		}
		w.Header().Add("Authorization", "Bearer "+token)
		utils.Ok(w, nil)
	} else {
		utils.Error(w, http.StatusNotAcceptable, "Incorrect password field")
	}
}

// TODO add validation
type signUpForm struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required,min=6"`
}

// SignUpHandler creates new user and returns jwt on success
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var form signUpForm
	err := utils.UnmarshalForm(r, &form)
	if err == utils.ErrIncorrectBody {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	} else if err == utils.ErrIncorectFormFields {
		utils.Error(w, http.StatusNotAcceptable, err.Error())
		return
	} else if err != nil {
		log.Panic("Assert: unexpected error: " + err.Error())
	}

	if !database.LoginIsFree(form.Username) {
		utils.Error(w, http.StatusConflict, "Username is not free")
		return
	}

	user, err := database.NewUser(form.FirstName, form.LastName, form.Username, form.Password)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Cannot create new user")
		return
	}
	err = user.Save()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Cannot create new user")
		return
	}
	utils.Ok(w, nil)
}

// UserInfoHandler returns info about user
var UserInfoHandler = middlewares.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user_id")
	if id == nil {
		log.Panic(errors.New("Assert: no user_id in context"))
		return
	}

	user, err := database.UserByID(id.(string))
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "Incorrect user id")
		return
	}

	utils.Ok(w, map[string]string{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"username":   user.Login,
	})
})

type changePasswordForm struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

// ChangePasswordHandler changes user password
var ChangePasswordHandler = middlewares.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user_id")
	if id == nil {
		log.Panic(errors.New("Assert: no user_id in context"))
		return
	}

	var form changePasswordForm
	err := utils.UnmarshalForm(r, &form)
	if err == utils.ErrIncorrectBody {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	} else if err == utils.ErrIncorectFormFields {
		utils.Error(w, http.StatusNotAcceptable, err.Error())
		return
	} else if err != nil {
		log.Panic("Assert: unexpected error: " + err.Error())
	}

	user, err := database.UserByID(id.(string))
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "Incorrect user id")
		return
	}

	if user.CheckPassword(form.OldPassword) {
		err := user.SetPassword(form.NewPassword)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Cannot store password")
			return
		}
		err = user.Save()
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Cannot store password")
			return
		}
		utils.Ok(w, nil)
	} else {
		utils.Error(w, http.StatusNotAcceptable, "Incorrect current password")
	}
})
