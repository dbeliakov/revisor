package config

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

var (
	Debug          = false
	SecretKey      = "SECRET_KEY"
	ClientFilesDir = "./client"
	DatabaseFile   = "/database/revisor.db"
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
	updateBoolFromEnv(&Debug, "DEBUG")
	updateFromEnv(&SecretKey, "SECRET_KEY")
	updateFromEnv(&ClientFilesDir, "CLIENT_FILES_DIR")
	if Debug {
		DatabaseFile = "./revisor.db"
	}
	updateFromEnv(&DatabaseFile, "DATABASE_FILE")
}
