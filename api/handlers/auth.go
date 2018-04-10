package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reviewer/api/auth"
)

type loginForm struct {
	Username string
	Password string
}

// LoginHandler checks username and password and returns jwt on success
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New login request")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var form loginForm
	err = json.Unmarshal(body, &form)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// TODO load from database by username
	user, err := auth.NewUser("Dmitrii", "Beliakov", form.Username, "abacaba", "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user.CheckPassword(form.Password) {
		token, err := user.NewToken("1")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Authorization", "Bearer "+token)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("GOOD"))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
