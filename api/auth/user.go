package auth

import (
	"github.com/dbeliakov/revisor/api/store"
	"golang.org/x/crypto/bcrypt"
)

const (
	hashCost = 8
)

// NewUser creates new user object
func newUser(firstName, lastName, login, password string) (store.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	if err != nil {
		return store.User{}, err
	}
	return store.User{
		FirstName:    firstName,
		LastName:     lastName,
		Login:        login,
		PasswordHash: string(hash),
	}, nil
}

// CheckPassword compares password with PasswordHash
func checkPassword(user store.User, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) == nil
}

// SetPassword new password
func setPassword(user *store.User, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hash)
	return nil
}
