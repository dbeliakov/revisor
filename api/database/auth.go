package database

import (
	"reviewer/api/auth"
	"reviewer/api/config"

	"github.com/globalsign/mgo/bson"
)

const (
	userCollectionName = "users"
)

// User mongo db model
type User struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	auth.User
}

// NewUser returns new mongo db user object
func NewUser(firstName, lastName, login, password string) (User, error) {
	user, err := auth.NewUser(firstName, lastName, login, password)
	if err != nil {
		return User{}, err
	}
	return User{
		User: user,
	}, nil
}

// UserByLogin finds user in db by login
func UserByLogin(login string) (User, error) {
	s := session.Copy()
	defer s.Close()

	c := s.DB(config.MongoDatabase).C(userCollectionName)
	var user User
	err := c.Find(bson.M{"user.login": login}).One(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

// UserByID finds user in db by id
func UserByID(id string) (User, error) {
	s := session.Copy()
	defer s.Close()

	c := s.DB(config.MongoDatabase).C(userCollectionName)
	var user User
	err := c.FindId(bson.ObjectIdHex(id)).One(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

// LoginIsFree checks if login is free
func LoginIsFree(login string) bool {
	s := session.Copy()
	defer s.Close()
	c := s.DB(config.MongoDatabase).C(userCollectionName)
	n, err := c.Find(bson.M{"user.login": login}).Count()
	if err != nil || n != 0 {
		return false
	}
	return true
}

// Save user info in db
func (user *User) Save() error {
	s := session.Copy()
	defer s.Close()

	c := s.DB(config.MongoDatabase).C(userCollectionName)
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
