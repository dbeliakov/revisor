package store

import (
	"github.com/asdine/storm"
	"github.com/dbeliakov/revisor/api/config"
	"golang.org/x/xerrors"
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
		panic(xerrors.Errorf("Cannot open database: %w", err))
	}
	Auth = newAuthStore(db)
	Comments = newCommentsStore(db)
	Reviews = newReviewsStore(db)
}
