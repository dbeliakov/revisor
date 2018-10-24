package store

import (
	"encoding/json"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// Comment represents information about comment
type Comment struct {
	ID      string `json:"id"`
	Author  string `json:"author"`
	Created int64  `json:"created"`
	Text    string `json:"text"`
	ParetID string `json:"parent;omitempty"`
	LineID  string `json:"line_id"`
}

// CommentsStore provides access to comments module storage
type CommentsStore interface {
	CreateComment(reviewID string, comment *Comment) error
	FindCommentByID(reviewID, id string) (Comment, error)
	CheckExists(reviewID, commentID string) (bool, error)
}

var (
	commentsBucket = []byte("comments")
)

type commentsStoreImpl struct {
	db *bolt.DB
}

func newCommentsStore(db *bolt.DB) CommentsStore {
	store := commentsStoreImpl{db: db}
	createBuckets(store.db, [][]byte{commentsBucket})
	return store
}

func (s commentsStoreImpl) FindCommentByID(reviewID, id string) (comment Comment, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(commentsBucket).Bucket([]byte(reviewID))
		if b == nil {
			return errors.New("No comments bucket for review")
		}
		v := b.Get([]byte(id))
		if v == nil {
			return errors.New("No comment with such id")
		}
		err := json.Unmarshal(v, &comment)
		if err != nil {
			return errors.Wrap(err, "Cannot deserialize comment")
		}
		return nil
	})
	return
}

func (s commentsStoreImpl) CreateComment(reviewID string, comment *Comment) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.Bucket(commentsBucket).CreateBucketIfNotExists([]byte(reviewID))
		if b == nil || err != nil {
			return errors.New("No comments bucket for review")
		}
		id, _ := b.NextSequence()
		comment.ID = strconv.FormatUint(id, 10)
		bytes, err := json.Marshal(comment)
		if err != nil {
			return errors.Wrap(err, "Cannot serialize comment")
		}
		err = b.Put([]byte(comment.ID), bytes)
		if err != nil {
			return errors.Wrap(err, "Cannot put new comment")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s commentsStoreImpl) CheckExists(reviewID, commentID string) (res bool, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(commentsBucket).Bucket([]byte(reviewID))
		if b == nil {
			return errors.New("No comments bucket for review")
		}
		res = b.Get([]byte(commentID)) != nil
		return nil
	})
	return
}
