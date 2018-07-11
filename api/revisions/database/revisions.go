package database

import (
	"encoding/json"
	"errors"
	auth "reviewer/api/auth/database"
	"reviewer/api/config"
	. "reviewer/api/database"
	"reviewer/api/revisions/lib"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const (
	collectionName = "reviews"
)

func collection(session *mgo.Session) *mgo.Collection {
	return session.DB(config.MongoDatabase).C(collectionName)
}

// Review is database object of review
type Review struct {
	ID       bson.ObjectId     `bson:"_id,omitempty" json:"id"`
	File     lib.VersionedFile `json:"-"`
	Name     string            `json:"name"`
	Updated  time.Time         `json:"-"`
	Closed   bool              `json:"closed"`
	Accepted bool              `json:"accepted"`

	OwnerID     bson.ObjectId
	ReviewersID []bson.ObjectId
}

// Owner user object
func (review Review) Owner() (auth.User, error) {
	user, err := auth.UserByID(review.OwnerID.Hex())
	if err != nil {
		return user, err
	}
	return user, nil
}

// Reviewers user objects
func (review Review) Reviewers() ([]auth.User, error) {
	var result []auth.User
	for _, reviewerID := range review.ReviewersID {
		user, err := auth.UserByID(reviewerID.Hex())
		if err != nil {
			return nil, err
		}
		result = append(result, user)
	}
	return result, nil
}

// MarshalJSON returns json with basic information about review
func (review Review) MarshalJSON() ([]byte, error) {
	owner, err := review.Owner()
	if err != nil {
		return nil, err
	}
	reviewers, err := review.Reviewers()
	if err != nil {
		return nil, err
	}

	type Alias Review
	return json.Marshal(&struct {
		Owner          auth.User   `json:"owner"`
		Reviewers      []auth.User `json:"reviewers"`
		Updated        int64       `json:"updated"`
		RevisionsCount int         `json:"revisions_count"`
		*Alias
	}{
		Owner:          owner,
		Reviewers:      reviewers,
		Updated:        review.Updated.Unix(),
		RevisionsCount: review.File.RevisionsCount(),
		Alias:          (*Alias)(&review),
	})
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
		File:        lib.NewVersionedFile(filename, content),
		Name:        name,
		OwnerID:     bson.ObjectIdHex(owner),
		ReviewersID: toBsonIds(reviewers),
		Updated:     time.Now(),
		Closed:      false,
		Accepted:    false,
	}
}

// ReviewByID from database
func ReviewByID(id string) (Review, error) {
	if !bson.IsObjectIdHex(id) {
		return Review{}, errors.New("Invalid bson ObjectId string")
	}
	s := Session.Copy()
	defer s.Close()

	c := collection(s)
	var file Review
	err := c.FindId(bson.ObjectIdHex(id)).One(&file)
	return file, err
}

// ReviewsByConditions with conditions
func ReviewsByConditions(condition bson.M) ([]Review, error) {
	s := Session.Copy()
	defer s.Close()

	c := s.DB(config.MongoDatabase).C(collectionName)
	files := make([]Review, 0)
	err := c.Find(condition).Sort("-updated").All(&files)
	if err != nil {
		return files, err
	}
	return files, nil
}

// Save review file to database
func (review *Review) Save() error {
	s := Session.Copy()
	defer s.Close()

	c := s.DB(config.MongoDatabase).C(collectionName)
	if !review.ID.Valid() {
		review.ID = bson.NewObjectId()
		err := c.Insert(nil, review)
		return err
	}

	err := c.Update(bson.M{"_id": review.ID}, review)
	if err != nil {
		return err
	}
	return nil
}
