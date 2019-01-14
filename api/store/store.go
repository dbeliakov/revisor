package store

import (
	"time"

	"github.com/boltdb/bolt"
	"github.com/dbeliakov/revisor/api/config"
	"github.com/pkg/errors"
)

var (
	// Auth module storage
	Auth AuthStore
	// Comments module storage
	Comments CommentsStore
	// Reviews module storage
	Reviews ReviewsStore
)

func init() {
	db, err := bolt.Open(config.DatabaseFile, 0666, &bolt.Options{Timeout: 30 * time.Second})
	if err != nil {
		panic(errors.Wrap(err, "Cannot open database"))
	}

	Auth = newAuthStore(db)
	Comments = newCommentsStore(db)
	Reviews = newReviewsStore(db)
}

func createBuckets(db *bolt.DB, buckets [][]byte) {
	err := db.Update(func(tx *bolt.Tx) error {
		for _, b := range buckets {
			_, err := tx.CreateBucketIfNotExists(b)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
