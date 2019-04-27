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
	tx, err := s.db.Begin(true)
	if err != nil {
		return errors.Wrap(err, "Cannot start transaction")
	}
	defer tx.Rollback()

	var exists User
	err = tx.One("Login", user.Login, &exists)
	if err != nil && err != storm.ErrNotFound {
		return errors.Wrap(err, "Cannot load user from database")
	} else if err == nil {
		return ErrUserExists
	}

	err = tx.Save(&user)
	if err != nil {
		return errors.Wrap(err, "Cannot save user")
	}
	tx.Commit()
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
	foundLogins := make(map[string]bool)

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

	result := make([]User, 0, searchUsersCount)

	for _, user := range byLogin {
		foundLogins[user.Login] = true
		if user.Login == exclude {
			continue
		}
		result = append(result, user)
		if len(result) == searchUsersCount {
			return result, nil
		}
	}

	for _, user := range byFirstName {
		if _, ok := foundLogins[user.Login]; ok {
			continue
		}
		foundLogins[user.Login] = true
		if user.Login == exclude {
			continue
		}
		result = append(result, user)
		if len(result) == searchUsersCount {
			return result, nil
		}
	}

	for _, user := range byLastName {
		if _, ok := foundLogins[user.Login]; ok {
			continue
		}
		foundLogins[user.Login] = true
		if user.Login == exclude {
			continue
		}
		result = append(result, user)
		if len(result) == searchUsersCount {
			return result, nil
		}
	}

	return result, nil
}
