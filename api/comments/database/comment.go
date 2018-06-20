package database

import (
	"encoding/json"
	"time"

	auth "reviewer/api/auth/database"
	"reviewer/api/config"
	. "reviewer/api/database"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const (
	userCollectionName = "comments"
)

func collection(session *mgo.Session) *mgo.Collection {
	return session.DB(config.MongoDatabase).C(userCollectionName)
}

// Comment database struct
type Comment struct {
	ID       bson.ObjectId   `bson:"_id,omitempty" json:"id"`
	AuthorID bson.ObjectId   `json:"-"`
	Created  time.Time       `json:"-"`
	Text     string          `json:"text"`
	ChildIDs []bson.ObjectId `json:"-"`
	ReviewID bson.ObjectId   `json:"review_id"`
	LineID   string          `json:"line_id"`
	Root     bool            `json:"-"`
}

// NewComment creates new comment
func NewComment(author string, text string, reviewID string, lineID string, root bool) Comment {
	return Comment{
		AuthorID: bson.ObjectIdHex(author),
		Created:  time.Now(),
		Text:     text,
		ReviewID: bson.ObjectIdHex(reviewID),
		LineID:   lineID,
		Root:     root,
	}
}

// Author returns user object
func (comment Comment) Author() (auth.User, error) {
	author, err := auth.UserByID(comment.AuthorID.Hex())
	if err != nil {
		return author, err
	}
	return author, nil
}

// Childs objects
func (comment Comment) Childs() ([]Comment, error) {
	var result []Comment
	for _, id := range comment.ChildIDs {
		child, err := CommentByID(id.Hex())
		if err != nil {
			return nil, err
		}
		result = append(result, child)
	}
	return result, nil
}

// Save comment in database
func (comment *Comment) Save() error {
	s := Session.Copy()
	defer s.Close()

	c := collection(s)
	if !comment.ID.Valid() {
		comment.ID = bson.NewObjectId()
		err := c.Insert(nil, comment)
		return err
	}

	err := c.Update(bson.M{"_id": comment.ID}, comment)
	if err != nil {
		return err
	}
	return nil
}

// MarshalJSON returns json with information about comment
func (comment Comment) MarshalJSON() ([]byte, error) {
	childs, err := comment.Childs()
	if err != nil {
		return nil, err
	}
	author, err := comment.Author()
	if err != nil {
		return nil, err
	}

	type Alias Comment
	return json.Marshal(&struct {
		Author  auth.User `json:"author"`
		Created int64     `json:"created"`
		Childs  []Comment `json:"childs"`
		*Alias
	}{
		Author:  author,
		Created: comment.Created.Unix(),
		Childs:  childs,
		Alias:   (*Alias)(&comment),
	})
}

// RootCommentsForReview returns comments without
func RootCommentsForReview(reviewID string) ([]Comment, error) {
	s := Session.Copy()
	defer s.Close()

	c := collection(s)
	comments := make([]Comment, 0)
	err := c.Find(bson.M{"reviewid": bson.ObjectIdHex(reviewID), "root": true}).All(&comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// CommentByID find comment by id
func CommentByID(id string) (Comment, error) {
	s := Session.Copy()
	defer s.Close()

	c := collection(s)
	var comment Comment
	err := c.FindId(bson.ObjectIdHex(id)).One(&comment)
	if err != nil {
		return comment, err
	}
	return comment, nil
}
