package database

import (
	"errors"
	"log"
	"reviewer/api/auth/lib"
	"reviewer/api/config"
	. "reviewer/api/database"

	"github.com/sirupsen/logrus"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const (
	userCollectionName = "users"
)

// Make 'login' field unique in database
func init() {
	s := Session.Copy()
	defer s.Close()

	loginField := "user.login"

	c := collection(s)
	indexes, err := c.Indexes()
	if err != nil {
		logrus.Warnf("Cannot get indexes: $+v. It's a new database?", err)
		log.Print("Error while getting indexes: ", err)
		indexes = []mgo.Index{}
	}
	for _, index := range indexes {
		for _, key := range index.Key {
			if key == loginField {
				return
			}
		}
	}

	index := mgo.Index{
		Key:    []string{loginField},
		Unique: true,
	}
	err = c.EnsureIndex(index)
	if err != nil {
		logrus.Fatalf("Error while creating index for 'login' field: %+v", err)
	}
}

func collection(session *mgo.Session) *mgo.Collection {
	return session.DB(config.MongoDatabase).C(userCollectionName)
}

// User mongo db model
type User struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"id"`
	lib.User
}

// NewUser returns new mongo db user object
func NewUser(firstName, lastName, login, password string) (User, error) {
	user, err := lib.NewUser(firstName, lastName, login, password)
	if err != nil {
		return User{}, err
	}
	return User{
		User: user,
	}, nil
}

// UserByLogin finds user in db by login
func UserByLogin(login string) (User, error) {
	s := Session.Copy()
	defer s.Close()

	c := collection(s)
	var user User
	err := c.Find(bson.M{"user.login": login}).One(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

// UserByID finds user in db by id
func UserByID(id string) (User, error) {
	if !bson.IsObjectIdHex(id) {
		return User{}, errors.New("Invalid bson ObjectId string")
	}
	s := Session.Copy()
	defer s.Close()

	c := collection(s)
	var user User
	err := c.FindId(bson.ObjectIdHex(id)).One(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

// LoginIsFree checks if login is free
func LoginIsFree(login string) bool {
	s := Session.Copy()
	defer s.Close()

	c := collection(s)
	n, err := c.Find(bson.M{"user.login": login}).Count()
	if err != nil || n != 0 {
		return false
	}
	return true
}

// Save user info in database
func (user *User) Save() error {
	s := Session.Copy()
	defer s.Close()

	c := collection(s)
	if !user.ID.Valid() {
		user.ID = bson.NewObjectId()
		err := c.Insert(nil, user)
		return err
	}

	err := c.Update(bson.M{"_id": user.ID}, user)
	if err != nil {
		return err
	}
	return nil
}

// SearchUsers by login or name
// TODO use indexes and text search
func SearchUsers(query string) ([]User, error) {
	s := Session.Copy()
	defer s.Close()

	const COUNT = 5

	c := collection(s)
	var results []User
	var tmpResults []User
	err := c.Find(bson.M{"user.login": &bson.RegEx{Pattern: query, Options: "i"}}).All(&results)
	if err != nil {
		return results, err
	}
	results = append(results, tmpResults...)
	if len(results) >= COUNT {
		return results[:COUNT], nil
	}
	err = c.Find(bson.M{"user.lastname": &bson.RegEx{Pattern: query, Options: "i"}}).All(&tmpResults)
	if err != nil {
		return results, err
	}
	results = append(results, tmpResults...)
	if len(results) >= COUNT {
		return results[:COUNT], nil
	}
	err = c.Find(bson.M{"user.firstname": &bson.RegEx{Pattern: query, Options: "i"}}).All(&tmpResults)
	if err != nil {
		return results, err
	}
	results = append(results, tmpResults...)
	if len(results) >= COUNT {
		return results[:COUNT], nil
	}
	return results, nil
}
