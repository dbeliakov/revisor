package database

import (
	"reviewer/api/auth/lib"
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
