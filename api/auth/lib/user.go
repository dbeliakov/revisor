package lib

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	hashCost = 8
)

// User represents information about registered user
type User struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Login        string `json:"username"`
	PasswordHash string `json:"-"`
}

// NewUser creates new user object
func NewUser(firstName, lastName, login, password string) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	if err != nil {
		return User{}, err
	}
	return User{
		FirstName:    firstName,
		LastName:     lastName,
		Login:        login,
		PasswordHash: string(hash),
	}, nil
}

// CheckPassword compares password with PasswordHash
func (user User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) == nil
}

// SetPassword new password
func (user *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hash)
	return nil
}
