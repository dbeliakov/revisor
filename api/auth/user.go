package auth

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	hashCost = 8
)

// User represents information about registered user
type User struct {
	FirstName    string
	LastName     string
	Login        string
	PasswordHash string
	Email        string
}

// NewUser creates new user object
func NewUser(firstName, lastName, login, password, email string) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	if err != nil {
		return User{}, err
	}
	return User{
		FirstName:    firstName,
		LastName:     lastName,
		Login:        login,
		PasswordHash: string(hash),
		Email:        email,
	}, nil
}

// CheckPassword compares password with PasswordHash
func (user User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) == nil
}
