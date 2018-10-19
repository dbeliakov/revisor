package store

import (
	"encoding/json"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// User represents information about registered user
type User struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Login        string `json:"username"`
	PasswordHash string `json:"password_hash;omitempty"`
}

// AuthStore provides access to auth module storage
type AuthStore interface {
	CreateBuckets() error

	FindUserByLogin(login string) (User, error)
	CreateUser(user User) error
	UpdateUser(user User) error
	FindUsers(query string) ([]User, error)
}

var (
	// ErrUserExists returns in case of login is not free
	ErrUserExists = errors.New("User with such login already exists")

	usersBucket           = []byte("users")
	searchFirstNameBucket = []byte("users-first-name")
	searchLastNameBucket  = []byte("users-last-name")
	loginListSeparator    = ","
)

type storeImpl struct {
	db *bolt.DB
}

func newAuthStore(db *bolt.DB) AuthStore {
	return storeImpl{db: db}
}

func (s storeImpl) CreateBuckets() error {
	buckets := [][]byte{usersBucket, searchFirstNameBucket, searchLastNameBucket}
	err := s.db.Update(func(tx *bolt.Tx) error {
		for _, b := range buckets {
			_, err := tx.CreateBucketIfNotExists(b)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s storeImpl) FindUserByLogin(login string) (User, error) {
	var user User
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(usersBucket)
		err := json.Unmarshal(b.Get([]byte(login)), &user)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return user, errors.Wrap(err, "Cannot find user")
	}
	return user, nil
}

func (s storeImpl) CreateUser(user User) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(usersBucket)
		if b.Get([]byte(user.Login)) != nil {
			return ErrUserExists
		}
		err := putUser(tx, user)
		if err != nil {
			return errors.Wrap(err, "Cannot put user")
		}
		return nil
	})
	if err == ErrUserExists {
		return err
	} else if err != nil {
		return errors.Wrap(err, "Cannot create new user")
	}
	return nil
}

func (s storeImpl) UpdateUser(user User) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		err := putUser(tx, user)
		if err != nil {
			return errors.Wrap(err, "Cannot put user")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s storeImpl) FindUsers(query string) ([]User, error) {
	// TODO
	return nil, nil
}

func putUser(tx *bolt.Tx, user User) error {
	b := tx.Bucket(usersBucket)
	serialized, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, "Cannot serialize user")
	}
	err = b.Put([]byte(user.Login), serialized)
	if err != nil {
		return errors.Wrap(err, "Cannot create user")
	}

	updateBucket := func(b *bolt.Bucket, key string) error {
		list := b.Get([]byte(key))
		if list == nil {
			list = make([]byte, 0)
		}
		list = []byte(strings.Join(
			append(strings.Split(string(list), loginListSeparator), user.Login),
			loginListSeparator))
		err := b.Put([]byte(key), list)
		if err != nil {
			return err
		}
		return nil
	}

	err = updateBucket(tx.Bucket(searchFirstNameBucket), user.FirstName)
	if err != nil {
		return err
	}
	err = updateBucket(tx.Bucket(searchLastNameBucket), user.LastName)
	if err != nil {
		return err
	}
	return nil
}
