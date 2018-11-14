package notifications

import (
	"net/http"
	"reviewer/api/auth"
	"reviewer/api/store"
	"reviewer/api/utils"

	"github.com/sirupsen/logrus"
)

// LinkTelegram to user account
var LinkTelegram = auth.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	u, err := auth.UserFromRequest(r)
	if err != nil {
		logrus.Errorf("Error while getting user from request context: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("No authorized user for this request"))
		return
	}
	user, err := store.Auth.FindUserByLogin(u.Login)
	if err != nil {
		logrus.Errorf("Error while getting user from database: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("No authorized user for this request"))
		return
	}

	var form struct {
		Login string `json:"username" validate:"required"`
		ID    int    `json:"id" validate:"required"`
	}
	if err := utils.UnmarshalForm(w, r, &form); err != nil {
		return
	}

	user.TelegramID = form.ID
	user.TelegramLogin = form.Login

	err = store.Auth.UpdateUser(user)
	if err != nil {
		logrus.Errorf("Cannot save telegram credentials")
		utils.Error(w, utils.InternalErrorResponse("Cannot save telegram credentials"))
		return
	}

	utils.Ok(w, nil)
})

// UnlinkTelegram from user account
var UnlinkTelegram = auth.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	u, err := auth.UserFromRequest(r)
	if err != nil {
		logrus.Errorf("Error while getting user from request context: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("No authorized user for this request"))
		return
	}
	user, err := store.Auth.FindUserByLogin(u.Login)
	if err != nil {
		logrus.Errorf("Error while getting user from database: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("No authorized user for this request"))
		return
	}

	user.TelegramID = 0
	user.TelegramLogin = ""

	err = store.Auth.UpdateUser(user)
	if err != nil {
		logrus.Errorf("Cannot save telegram credentials")
		utils.Error(w, utils.InternalErrorResponse("Cannot save telegram credentials"))
		return
	}

	utils.Ok(w, nil)
})

// TelegramLogin of user; if exists, user linked telegram for notifications
var TelegramLogin = auth.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
	u, err := auth.UserFromRequest(r)
	if err != nil {
		logrus.Errorf("Error while getting user from request context: %+v", err)
		utils.Error(w, utils.InternalErrorResponse("No authorized user for this request"))
		return
	}

	utils.Ok(w, map[string]string{
		"login": u.TelegramLogin,
	})
})
