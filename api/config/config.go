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
	DatabaseFile   = "./revisor.db"
	TgAPIKey       = "secret_key"
	BaseURL        = "http://localhost:8080"
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
	updateFromEnv(&DatabaseFile, "DATABASE_FILE")
	updateFromEnv(&TgAPIKey, "TG_API_KEY")
	updateFromEnv(&BaseURL, "BASE_URL")
	TgAPIKey = "741718821:AAGX3Qt2ICZw3RacIiTree-3r97qN_zM-5g" // TODO
}
