package middlewares

import (
	"context"
	"net/http"
	"reviewer/api/auth"
	"reviewer/api/utils"
)

// AuthRequired checks jwt token
func AuthRequired(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authString := r.Header.Get("Authorization")
		bearerLength := len("Bearer ")
		if len(authString) < bearerLength {
			utils.Unauthorized(w)
			return
		}
		tokenString := authString[bearerLength:]
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			utils.Unauthorized(w)
			return
		}
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user_id", claims["id"])))
	}
}
