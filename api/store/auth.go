package store

import (
	"bytes"
	"encoding/json"

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
	FindUserByLogin(login string) (User, error)
	CreateUser(user User) error
	UpdateUser(user User) error
	FindUsers(query string) ([]User, error)
	CheckExists(login string) (bool, error)
}

var (
	// ErrUserExists returns in case of login is not free
	ErrUserExists = errors.New("User with such login already exists")

	usersBucket           = []byte("users")
	searchFirstNameBucket = []byte("users-first-name")
	searchLastNameBucket  = []byte("users-last-name")
	loginListSeparator    = []byte(",")
	searchUsersCount      = 5
)

type authStoreImpl struct {
	db *bolt.DB
}

func newAuthStore(db *bolt.DB) AuthStore {
	store := authStoreImpl{db: db}
	createBuckets(store.db, [][]byte{usersBucket, searchFirstNameBucket, searchLastNameBucket})
	return store
}

func (s authStoreImpl) FindUserByLogin(login string) (User, error) {
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

func (s authStoreImpl) CheckExists(login string) (res bool, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(usersBucket)
		res = b.Get([]byte(login)) != nil
		return nil
	})
	return
}

func (s authStoreImpl) CreateUser(user User) error {
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

// TODO only for password changing now
func (s authStoreImpl) UpdateUser(user User) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(usersBucket)
		serialized, err := json.Marshal(user)
		if err != nil {
			return errors.Wrap(err, "Cannot serialize user")
		}
		err = b.Put([]byte(user.Login), serialized)
		if err != nil {
			return errors.Wrap(err, "Cannot update user")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s authStoreImpl) FindUsers(query string) ([]User, error) {
	result := make([]User, 0)
	err := s.db.View(func(tx *bolt.Tx) error {
		users := tx.Bucket(usersBucket)
		c := users.Cursor()
		prefix := []byte(query)
		for k, v := c.Seek(prefix); k != nil &&
			bytes.HasPrefix(k, prefix) && len(result) < searchUsersCount; k, v = c.Next() {
			result = append(result, User{})
			err := json.Unmarshal(v, &result[len(result)-1])
			if err != nil {
				return errors.Wrap(err, "Cannot deserialize user")
			}
		}

		addResults := func(c *bolt.Cursor) error {
			for k, v := c.Seek(prefix); k != nil &&
				bytes.HasPrefix(k, prefix) && len(result) < searchUsersCount; k, v = c.Next() {
				logins := bytes.Split(v, []byte(loginListSeparator))
				for _, login := range logins {
					result = append(result, User{})
					err := json.Unmarshal(users.Get(login), &result[len(result)-1])
					if err != nil {
						return errors.Wrap(err, "Cannot deserialize user")
					}
				}
			}
			return nil
		}
		err := addResults(tx.Bucket(searchLastNameBucket).Cursor())
		if err != nil {
			return err
		}
		err = addResults(tx.Bucket(searchFirstNameBucket).Cursor())
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
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
		splitted := make([][]byte, 0)
		if list != nil && len(list) > 0 {
			splitted = bytes.Split(list, loginListSeparator)
		}
		list = bytes.Join(
			append(splitted, []byte(user.Login)),
			loginListSeparator)
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
