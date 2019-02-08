package store

import (
	"github.com/asdine/storm"
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

// InitStore and open database
func InitStore() {
	db, err := storm.Open(config.DatabaseFile)
	if err != nil {
		panic(errors.Wrap(err, "Cannot open database"))
	}
	Auth = newAuthStore(db)
	Comments = newCommentsStore(db)
	Reviews = newReviewsStore(db)
}
