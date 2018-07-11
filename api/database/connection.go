package database

import (
	"reviewer/api/config"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/globalsign/mgo"
)

var (
	// Session to be shared between components
	Session *mgo.Session
)

func init() {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{config.MongoAddr},
		Timeout:  10 * time.Second,
		Database: config.MongoDatabase,
	}
	var err error
	Session, err = mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		logrus.Panic("Cannot connect to database: %+v", err)
	}
}
