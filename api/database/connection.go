package database

import (
	"log"
	"reviewer/api/config"
	"time"

	"github.com/globalsign/mgo"
)

var (
	// Session to be shared between components
	Session *mgo.Session
)

func init() {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{config.MongoAddr},
		Timeout:  60 * time.Second,
		Database: config.MongoDatabase,
	}
	var err error
	Session, err = mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Panic(err)
	}
}
