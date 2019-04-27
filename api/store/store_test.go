package store

import (
	"os"

	"github.com/asdine/storm"
	"github.com/pkg/errors"
)

const (
	testDatabase = "test_database.db"
)

func initTestDatabase() {
	db, err := storm.Open(testDatabase)
	if err != nil {
		panic(errors.Wrap(err, "Cannot open test database"))
	}
	Auth = newAuthStore(db)
	Comments = newCommentsStore(db)
	Reviews = newReviewsStore(db)
}

func removeTestDatabase() {
	err := os.Remove(testDatabase)
	if err != nil {
		panic(errors.Wrap(err, "Cannot remove test database"))
	}
}
