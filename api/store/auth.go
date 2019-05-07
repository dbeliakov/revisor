package store

import (
	"github.com/asdine/storm"
	"golang.org/x/xerrors"
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
	ErrUserExists = xerrors.New("User with such login already exists")

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
		return user, xerrors.Errorf("Cannot find user: %w", err)
	}
	return user, nil
}

func (s authStoreImpl) CreateUser(user User) error {
	tx, err := s.db.Begin(true)
	if err != nil {
		return xerrors.Errorf("Cannot start transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var exists User
	err = tx.One("Login", user.Login, &exists)
	if err != nil && err != storm.ErrNotFound {
		return xerrors.Errorf("Cannot load user from database: %w", err)
	} else if err == nil {
		return xerrors.Errorf("user already exists: %w", ErrUserExists)
	}

	err = tx.Save(&user)
	if err != nil {
		return xerrors.Errorf("Cannot save user: %w")
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (s authStoreImpl) UpdateUser(user User) error {
	err := s.db.Update(&user)
	if err != nil {
		return xerrors.Errorf("Cannot update user: %w", err)
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
		return nil, xerrors.Errorf("Cannot select users by login: %w", err)
	}
	err = s.db.Prefix("FirstName", query, &byFirstName, storm.Limit(searchUsersCount+len(exclude)))
	if err != storm.ErrNotFound && err != nil {
		return nil, xerrors.Errorf("Cannot select users by first name: %w", err)
	}
	err = s.db.Prefix("LastName", query, &byLastName, storm.Limit(searchUsersCount+len(exclude)))
	if err != storm.ErrNotFound && err != nil {
		return nil, xerrors.Errorf("Cannot select users by last name: %w", err)
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
