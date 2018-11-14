package auth

import (
	"context"
	"errors"
	"net/http"
	"reviewer/api/store"
	"reviewer/api/utils"
	"time"

	"github.com/sirupsen/logrus"
)

// authMiddlewareKey type for request context
type authMiddlewareKey int

const (
	// keyUserID key for request context
	keyUser authMiddlewareKey = iota
)

// AuthRequired checks jwt token and sets user_id value to request context
func AuthRequired(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authString := r.Header.Get("Authorization")
		bearerLength := len("Bearer ")
		if len(authString) < bearerLength {
			logrus.Warnf("Incorrect length of \"Authorization\" header: %d", len(authString))
			utils.Unauthorized(w)
			return
		}
		tokenString := authString[bearerLength:]
		claims, err := validateToken(tokenString)
		if err != nil {
			logrus.Warnf("Cannot validate JWT-token: %+v", err)
			utils.Unauthorized(w)
			return
		}
		user := userFromToken(claims)
		if claims.VerifyExpiresAt(time.Now().Add(5*time.Hour*24).Unix(), false) {
			token, err := NewToken(user)
			if err != nil {
				logrus.Errorf("Cannot create new token for user: %s, error: %+v", user.Login, err)
			} else {
				w.Header().Add("Authorization", "Bearer "+token)
			}
		}
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), keyUser, user)))
	}
}

// UserFromRequest extracts user with UserID, specified in request context
func UserFromRequest(r *http.Request) (store.User, error) {
	u := r.Context().Value(keyUser)
	if u == nil {
		return store.User{}, errors.New("No \"keyUser\" value in request context")
	}
	// Load additional user information from database
	user, err := store.Auth.FindUserByLogin(u.(store.User).Login)
	if err != nil {
		return store.User{}, err
	}
	return user, nil
}
