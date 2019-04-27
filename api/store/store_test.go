package store

import (
	"os"

	"github.com/asdine/storm"
	"golang.org/x/xerrors"
)

const (
	testDatabase = "test_database.db"
)

func initTestDatabase() {
	db, err := storm.Open(testDatabase)
	if err != nil {
		panic(xerrors.Errorf("Cannot open test database: %w", err))
	}
	Auth = newAuthStore(db)
	Comments = newCommentsStore(db)
	Reviews = newReviewsStore(db)
}

func removeTestDatabase() {
	err := os.Remove(testDatabase)
	if err != nil {
		panic(xerrors.Errorf("Cannot remove test database: %w", err))
	}
}
