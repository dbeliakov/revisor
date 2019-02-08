package store

import (
	"github.com/asdine/storm"
	"github.com/pkg/errors"
)

// User represents information about registered user
type User struct {
	FirstName    string `storm:"index"`
	LastName     string `storm:"index"`
	Login        string `storm:"id"`
	PasswordHash string
}

// AuthStore provides access to auth module storage
type AuthStore interface {
	FindUserByLogin(login string) (User, error)
	CreateUser(user User) error
	UpdateUser(user User) error
	FindUsers(query string, exclude string) ([]User, error)
}

var (
	// ErrUserExists returns in case of login is not free
	ErrUserExists = errors.New("User with such login already exists")

	searchUsersCount = 5
)

type authStoreImpl struct {
	db *storm.DB
}

func newAuthStore(db *storm.DB) AuthStore {
	return authStoreImpl{db: db}
}

func (s authStoreImpl) FindUserByLogin(login string) (User, error) {
	var user User
	err := s.db.One("Login", login, &user)
	if err != nil {
		return user, errors.Wrap(err, "Cannot find user")
	}
	return user, nil
}

func (s authStoreImpl) CreateUser(user User) error {
	err := s.db.Save(&user)
	if err == storm.ErrAlreadyExists {
		return ErrUserExists
	}
	if err != nil {
		return errors.Wrap(err, "Cannot save user")
	}
	return nil
}

func (s authStoreImpl) UpdateUser(user User) error {
	err := s.db.Update(&user)
	if err != nil {
		return errors.Wrap(err, "Cannot update user")
	}
	return nil
}

func (s authStoreImpl) FindUsers(query string, exclude string) ([]User, error) {
	result := make([]User, 0, searchUsersCount)

	byLogin := make([]User, 0, searchUsersCount+len(exclude))
	byFirstName := make([]User, 0, searchUsersCount+len(exclude))
	byLastName := make([]User, 0, searchUsersCount+len(exclude))
	err := s.db.Prefix("Login", query, &byLogin, storm.Limit(searchUsersCount+len(exclude)))
	if err != storm.ErrNotFound && err != nil {
		return nil, errors.Wrap(err, "Cannot select users by login")
	}
	err = s.db.Prefix("FirstName", query, &byFirstName, storm.Limit(searchUsersCount+len(exclude)))
	if err != storm.ErrNotFound && err != nil {
		return nil, errors.Wrap(err, "Cannot select users by first name")
	}
	err = s.db.Prefix("LastName", query, &byLastName, storm.Limit(searchUsersCount+len(exclude)))
	if err != storm.ErrNotFound && err != nil {
		return nil, errors.Wrap(err, "Cannot select users by last name")
	}

	byLogin = append(byLogin, byFirstName...)
	byLogin = append(byLogin, byLastName...)

	for _, res := range byLogin {
		if res.Login != exclude {
			result = append(result, res)
			if len(result) == searchUsersCount {
				break
			}
		}
	}
	return result, nil
}
