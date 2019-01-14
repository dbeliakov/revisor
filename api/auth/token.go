package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/dbeliakov/revisor/api/config"
	"github.com/dbeliakov/revisor/api/store"
	"github.com/dgrijalva/jwt-go"
)

var (
	signingKey = []byte(config.SecretKey)
)

const (
	firstNameKey = "first_name"
	lastNameKey  = "last_name"
	loginKey     = "login"
	expiredTTL   = 7 * 24 * time.Hour
	refreshTTL   = 5 * 24 * time.Hour
)

// NewToken generates new JWT token for user
func newToken(user store.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		firstNameKey: user.FirstName,
		lastNameKey:  user.LastName,
		loginKey:     user.Login,
		"exp":        time.Now().Add(expiredTTL).Unix(),
	})
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateToken from user
func validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return signingKey, nil
	})
	if err != nil {
		return jwt.MapClaims{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return jwt.MapClaims{}, errors.New("Token is not valid")
}

// UserFromToken builds user object using information in token claims
func userFromToken(claims jwt.MapClaims) store.User {
	return store.User{
		FirstName: claims[firstNameKey].(string),
		LastName:  claims[lastNameKey].(string),
		Login:     claims[loginKey].(string),
	}
}
