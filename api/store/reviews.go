package store

import (
	"github.com/asdine/storm"
	"golang.org/x/xerrors"
)

// Review represents information about review
type Review struct {
	ID       int    `storm:"id,increment"`
	File     []byte // TODO store id of versioned file
	Name     string
	Updated  int64
	Closed   bool
	Accepted bool
	Owner    string `storm:"index"`
	// TODO create custom index for reviewers
	Reviewers []string
}

// ReviewsStore provides access to comments module storage
type ReviewsStore interface {
	CreateReview(review *Review) error
	FindReviewByID(id int) (Review, error)
	FindReviewsByOwner(owner string) ([]Review, error)
	FindReviewsByReviewer(reviewer string) ([]Review, error)
	UpdateReview(review *Review) error
}

type reviewsStoreImpl struct {
	db *storm.DB
}

func newReviewsStore(db *storm.DB) ReviewsStore {
	return reviewsStoreImpl{db: db}
}

func (s reviewsStoreImpl) CreateReview(review *Review) error {
	err := s.db.Save(review)
	if err != nil {
		return xerrors.Errorf("Cannot save review: %w", err)
	}
	return nil
}

func (s reviewsStoreImpl) UpdateReview(review *Review) error {
	err := s.db.Update(review)
	if err != nil {
		return xerrors.Errorf("Cannot update review: %w", err)
	}
	return nil
}

func (s reviewsStoreImpl) FindReviewByID(id int) (Review, error) {
	var review Review
	err := s.db.One("ID", id, &review)
	if err != nil {
		return review, xerrors.Errorf("Cannot find user by ID: %w", err)
	}
	return review, nil
}

func (s reviewsStoreImpl) FindReviewsByOwner(owner string) ([]Review, error) {
	reviews := make([]Review, 0)
	err := s.db.Find("Owner", owner, &reviews)
	if err == storm.ErrNotFound {
		return reviews, nil
	}
	if err != nil {
		return reviews, xerrors.Errorf("Cannot find users by Owner: %w", err)
	}
	return reviews, nil
}

func (s reviewsStoreImpl) FindReviewsByReviewer(reviewer string) ([]Review, error) {
	reviews := make([]Review, 0)
	err := s.db.Select().Each(new(Review), func(v interface{}) error {
		review := v.(*Review)
		for _, rv := range review.Reviewers {
			if rv == reviewer {
				reviews = append(reviews, *review)
				break
			}
		}
		return nil
	})
	if err != nil {
		return nil, xerrors.Errorf("Cannot find reviews by reviewer: %w", err)
	}
	return reviews, err
}
