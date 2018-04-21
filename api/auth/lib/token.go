package lib

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	signingKey = "test_signing_key"
)

// NewToken generates new JWT token
func (user *User) NewToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         id,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"login":      user.Login,
		"exp":        time.Now().Add(time.Hour * 48).Unix(),
	})
	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateToken from user
func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return jwt.MapClaims{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return jwt.MapClaims{}, errors.New("Token is not valid")
}
