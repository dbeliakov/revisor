package database

import (
	"log"
	"reviewer/api/revisions"
	"time"

	"github.com/globalsign/mgo/bson"

	"github.com/globalsign/mgo"
)

var (
	session *mgo.Session
)

// TODO: config
const (
	mongoAddr     = "localhost:27017"
	mongoDatabase = "reviewer"
)

const (
	collectionName = "versioned_files"
)

func init() {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{mongoAddr},
		Timeout:  60 * time.Second,
		Database: mongoDatabase,
	}
	var err error
	session, err = mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Panic(err)
	}
}

// VersionedFile is database object of versioned file
type VersionedFile struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	revisions.VersionedFile
}

// NewVersionedFile creates versioned file which can be stored in db
func NewVersionedFile(filename string, content []string) VersionedFile {
	return VersionedFile{
		VersionedFile: revisions.NewVersionedFile(filename, content),
	}
}

// LoadVersionedFile from database
func LoadVersionedFile(id string) (VersionedFile, error) {
	s := session.Copy()
	defer s.Close()

	c := s.DB(mongoDatabase).C(collectionName)
	var file VersionedFile
	err := c.FindId(bson.ObjectIdHex(id)).One(&file)
	return file, err
}

// Save versioned file to database
func (file *VersionedFile) Save() error {
	s := session.Copy()
	defer s.Close()

	c := s.DB(mongoDatabase).C(collectionName)
	if !file.ID.Valid() {
		file.ID = bson.NewObjectId()
		err := c.Insert(nil, file)
		return err
	}

	err := c.Update(bson.M{"_id": file.ID}, file)
	return err
}
