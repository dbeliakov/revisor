package database

import (
	"reviewer/api/config"
	"reviewer/api/revisions"
	"time"

	"github.com/globalsign/mgo/bson"
)

const (
	collectionName = "reviews"
)

// Review is database object of review
type Review struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	File      revisions.VersionedFile
	Name      string
	Owner     bson.ObjectId
	Reviewers []bson.ObjectId
	Updated   time.Time
	Closed    bool
	Accepted  bool
}

func toBsonIds(ids []string) []bson.ObjectId {
	var result []bson.ObjectId
	for _, id := range ids {
		result = append(result, bson.ObjectIdHex(id))
	}
	return result
}

// NewReview creates versioned file which can be stored in db
func NewReview(filename string, content []string, name string, owner string, reviewers []string) Review {
	return Review{
		File:      revisions.NewVersionedFile(filename, content),
		Name:      name,
		Owner:     bson.ObjectIdHex(owner),
		Reviewers: toBsonIds(reviewers),
		Updated:   time.Now(),
		Closed:    false,
		Accepted:  false,
	}
}

// LoadReview from database
func LoadReview(id string) (Review, error) {
	s := session.Copy()
	defer s.Close()

	c := s.DB(config.MongoDatabase).C(collectionName)
	var file Review
	err := c.FindId(bson.ObjectIdHex(id)).One(&file)
	return file, err
}

// FindReviews with conditions
func FindReviews(condition bson.M) ([]Review, error) {
	s := session.Copy()
	defer s.Close()

	c := s.DB(config.MongoDatabase).C(collectionName)
	var files []Review
	err := c.Find(condition).All(&files)
	if err != nil {
		return files, err
	}
	return files, nil
}

// Save versioned file to database
func (file *Review) Save() error {
	s := session.Copy()
	defer s.Close()

	c := s.DB(config.MongoDatabase).C(collectionName)
	if !file.ID.Valid() {
		file.ID = bson.NewObjectId()
		err := c.Insert(nil, file)
		return err
	}

	err := c.Update(bson.M{"_id": file.ID}, file)
	if err != nil {
		return err
	}
	return nil
}
