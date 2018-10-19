package store

import (
	"reviewer/api/config"
	"time"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

type store struct {
	Auth AuthStore
	// TODO
}

var (
	// Store is an object for accesing database
	Store = newStore()
)

func newStore() store {
	db, err := bolt.Open(config.DatabaseFile, 0666, &bolt.Options{Timeout: 30 * time.Second})
	if err != nil {
		panic(errors.Wrap(err, "Cannot open database"))
	}
	result := store{
		Auth: newAuthStore(db),
	}
	err = result.Auth.CreateBuckets()
	if err != nil {
		panic(errors.Wrap(err, "Cannot open database"))
	}
	return result
}
