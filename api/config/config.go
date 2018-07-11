package config

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

var (
	MongoAddr      = "localhost:27017"
	MongoDatabase  = "revisor"
	Debug          = false
	SecretKey      = "SECRET_KEY"
	ClientFilesDir = "./client"
)

func updateFromEnv(val *string, key string) {
	v := os.Getenv(key)
	if len(v) > 0 {
		*val = v
	}
}

func updateBoolFromEnv(val *bool, key string) {
	v := os.Getenv(key)
	if len(v) > 0 {
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			logrus.Warnf("Cannot parse env variable \"%s\": %+v", key, err)
			return
		}
		*val = parsed
	}
}

func init() {
	updateFromEnv(&MongoAddr, "MONGO_ADDR")
	updateFromEnv(&MongoDatabase, "MONGO_DATABASE")
	updateBoolFromEnv(&Debug, "DEBUG")
	updateFromEnv(&SecretKey, "SECRET_KEY")
	updateFromEnv(&ClientFilesDir, "CLIENT_FILES_DIR")
}
