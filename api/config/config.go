package config

import (
	"flag"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

var (
	// Debug start server in development mode
	Debug = false
	// SecretKey for JWT tokens
	SecretKey = "SECRET_KEY"
	// ClientFilesDir with static client files
	ClientFilesDir = "./client"
	// DatabaseFile path to database
	DatabaseFile = "/database/revisor.db"
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
	if flag.Lookup("test.v") != nil {
		Debug = true
	}
	updateBoolFromEnv(&Debug, "DEBUG")
	updateFromEnv(&SecretKey, "SECRET_KEY")
	updateFromEnv(&ClientFilesDir, "CLIENT_FILES_DIR")
	if Debug {
		DatabaseFile = "./revisor.db"
	}
	updateFromEnv(&DatabaseFile, "DATABASE_FILE")
}
