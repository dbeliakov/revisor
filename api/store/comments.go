package store

import (
	"strconv"

	"github.com/asdine/storm"
	"github.com/pkg/errors"
)

// Comment represents information about comment
type Comment struct {
	ID       int `storm:"id,increment"`
	Author   string
	Created  int64
	Text     string
	ParentID int
	LineID   string
}

// CommentsStore provides access to comments module storage
type CommentsStore interface {
	CreateComment(reviewID int, comment *Comment) error
	FindCommentByID(reviewID, id int) (Comment, error)
	CheckExists(reviewID, commentID int) (bool, error)
	CommentsForReview(reviewID int) ([]Comment, error)
}

type commentsStoreImpl struct {
	db *storm.DB
}

func newCommentsStore(db *storm.DB) CommentsStore {
	return commentsStoreImpl{db: db}
}

func (s commentsStoreImpl) node(reviewID int) storm.Node {
	return s.db.From(strconv.Itoa(reviewID))
}

func (s commentsStoreImpl) FindCommentByID(reviewID, id int) (Comment, error) {
	var comment Comment
	err := s.node(reviewID).One("ID", id, &comment)
	if err != nil {
		return comment, errors.Wrap(err, "Cannot find comment by ID")
	}
	return comment, nil
}

func (s commentsStoreImpl) CreateComment(reviewID int, comment *Comment) error {
	err := s.node(reviewID).Save(comment)
	if err != nil {
		return errors.Wrap(err, "Cannot save comment")
	}
	return nil
}

func (s commentsStoreImpl) CheckExists(reviewID, id int) (bool, error) {
	var comment Comment
	err := s.node(reviewID).One("ID", id, &comment)
	if err == storm.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrap(err, "Cannot find comment by ID")
	}
	return true, nil
}

func (s commentsStoreImpl) CommentsForReview(reviewID int) ([]Comment, error) {
	comments := make([]Comment, 0)
	err := s.node(reviewID).All(&comments)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot load all comments for review")
	}
	return comments, nil
}
