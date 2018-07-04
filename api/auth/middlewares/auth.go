package middlewares

import (
	"context"
	"errors"
	"net/http"
	"reviewer/api/auth/database"
	"reviewer/api/auth/lib"
	"reviewer/api/utils"

	"github.com/sirupsen/logrus"
)

// authMiddlewareKey type for request context
type authMiddlewareKey int

const (
	// keyUserID key for request context
	keyUserID authMiddlewareKey = iota
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
		claims, err := lib.ValidateToken(tokenString)
		if err != nil {
			logrus.Warnf("Cannot validate JWT-token: %+v", err)
			utils.Unauthorized(w)
			return
		}
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), keyUserID, claims["id"])))
	}
}

// UserFromRequest extracts user with UserID, specified in request context
func UserFromRequest(r *http.Request) (*database.User, error) {
	userID := r.Context().Value(keyUserID)
	if userID == nil {
		return nil, errors.New("No \"keyUserID\" value in request context")
	}
	user, err := database.UserByID(userID.(string))
	if err != nil {
		return nil, err
	}
	return &user, nil
}
