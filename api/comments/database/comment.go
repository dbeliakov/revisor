package database

import (
	"time"

	auth "reviewer/api/auth/database"
	"reviewer/api/config"
	. "reviewer/api/database"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const (
	userCollectionName = "users"
)

func collection(session *mgo.Session) *mgo.Collection {
	return session.DB(config.MongoDatabase).C(userCollectionName)
}

// Comment database struct
type Comment struct {
	ID       bson.ObjectId   `bson:"_id,omitempty" json:"id"`
	AuthorID bson.ObjectId   `json:"-"`
	Created  time.Time       `json:"created"`
	Text     string          `json:"text"`
	Childs   []bson.ObjectId `json:"childs"`
	IsOpen   bool            `json:"is_open"`
}

// NewComment creates new comment
func NewComment(author string, text string) Comment {
	return Comment{
		AuthorID: bson.ObjectIdHex(author),
		Created:  time.Now(),
		Text:     text,
		IsOpen:   true,
	}
}

// Author returns user object
func (c Comment) Author() (auth.User, error) {
	author, err := auth.UserByID(c.AuthorID.Hex())
	if err != nil {
		return author, err
	}
	return author, nil
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
