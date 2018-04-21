package middlewares

import (
	"context"
	"net/http"
	"reviewer/api/auth/lib"
	"reviewer/api/utils"
)

// AuthRequired checks jwt token and sets user_id value to request context
func AuthRequired(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authString := r.Header.Get("Authorization")
		bearerLength := len("Bearer ")
		if len(authString) < bearerLength {
			utils.Unauthorized(w)
			return
		}
		tokenString := authString[bearerLength:]
		claims, err := lib.ValidateToken(tokenString)
		if err != nil {
			utils.Unauthorized(w)
			return
		}
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user_id", claims["id"])))
	}
}
