package auth

import (
	"errors"
	"log"
	"net/http"
	"reviewer/api/auth/database"
	"reviewer/api/auth/middlewares"
	"reviewer/api/utils"
)

// LoginHandler checks username and password and returns jwt on success
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var form struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	if err := utils.UnmarshalForm(w, r, &form); err != nil {
		return
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

	utils.Ok(w, user)
})

// ChangePasswordHandler changes user password
var ChangePasswordHandler = middlewares.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user_id")
	if id == nil {
		log.Panic(errors.New("Assert: no user_id in context"))
		return
	}

	var form struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=6"`
	}
	if err := utils.UnmarshalForm(w, r, &form); err != nil {
		return
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
