package store

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"github.com/valyala/fastjson"
)

// Review represents information about review
type Review struct {
	ID        string   `json:"id"`
	File      []byte   `json:"versioned_file"`
	Name      string   `json:"name"`
	Updated   int64    `json:"updated"`
	Closed    bool     `json:"closed"`
	Accepted  bool     `json:"accepted"`
	Owner     string   `json:"owner"`
	Reviewers []string `json:"reviewers"`
}

// ReviewsStore provides access to comments module storage
type ReviewsStore interface {
	CreateReview(review *Review) error
	FindReviewByID(id string) (Review, error)
	FindReviewsByOwner(owner string) ([]Review, error)
	FindReviewsByReviewer(reviewer string) ([]Review, error)
	UpdateReview(review Review) error
}

var (
	reviewsBucket = []byte("reviews")
)

type reviewsStoreImpl struct {
	db *bolt.DB
}

func newReviewsStore(db *bolt.DB) ReviewsStore {
	store := reviewsStoreImpl{db: db}
	createBuckets(store.db, [][]byte{reviewsBucket})
	return store
}

func (r reviewsStoreImpl) CreateReview(review *Review) error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(reviewsBucket)
		id, _ := b.NextSequence()
		review.ID = strconv.FormatUint(id, 10)
		bytes, err := json.Marshal(review)
		if err != nil {
			return errors.Wrap(err, "Cannot serialize review")
		}
		err = b.Put([]byte(review.ID), bytes)
		if err != nil {
			return errors.Wrap(err, "Cannot put new review")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r reviewsStoreImpl) UpdateReview(review Review) error {
	err := r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(reviewsBucket)
		bytes, err := json.Marshal(review)
		if err != nil {
			return errors.Wrap(err, "Cannot serialize review")
		}
		err = b.Put([]byte(review.ID), bytes)
		if err != nil {
			return errors.Wrap(err, "Cannot put new review")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r reviewsStoreImpl) FindReviewByID(id string) (review Review, err error) {
	err = r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(reviewsBucket)
		v := b.Get([]byte(id))
		if v == nil {
			return errors.New("No review with such id")
		}
		err := json.Unmarshal(v, &review)
		if err != nil {
			return err
		}
		return nil
	})
	return
}

func (r reviewsStoreImpl) FindReviewsByOwner(owner string) ([]Review, error) {
	result := make([]Review, 0)
	ownerB := []byte(owner)
	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(reviewsBucket)
		c := b.Cursor()
		var p fastjson.Parser
		for k, v := c.First(); k != nil; k, v = c.Next() {
			parsed, err := p.ParseBytes(v)
			if err != nil {
				return errors.Wrap(err, "Cannot parse review")
			}
			reviewOwner := parsed.GetStringBytes("owner")
			if bytes.Equal(reviewOwner, ownerB) {
				result = append(result, Review{})
				err := json.Unmarshal(v, &result[len(result)-1])
				if err != nil {
					return errors.Wrap(err, "Cannot deserialize review")
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r reviewsStoreImpl) FindReviewsByReviewer(reviewer string) ([]Review, error) {
	result := make([]Review, 0)
	reviewerB := []byte(reviewer)
	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(reviewsBucket)
		c := b.Cursor()
		var p fastjson.Parser
		for k, v := c.First(); k != nil; k, v = c.Next() {
			parsed, err := p.ParseBytes(v)
			if err != nil {
				return errors.Wrap(err, "Cannot parse review")
			}
			reviewers := parsed.GetArray("reviewers")
			for _, r := range reviewers {
				if bytes.Equal(r.GetStringBytes(), reviewerB) {
					result = append(result, Review{})
					err := json.Unmarshal(v, &result[len(result)-1])
					if err != nil {
						return errors.Wrap(err, "Cannot deserialize review")
					}
					break
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
