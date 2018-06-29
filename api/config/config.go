package config

import "os"

var (
	MongoAddr     = "localhost:27017"
	MongoDatabase = "revisor"
	Debug         = false
)

func updateFromEnv(val *string, key string) {
	v := os.Getenv(key)
	if len(v) > 0 {
		*val = v
	}
}

func init() {
	updateFromEnv(&MongoAddr, "MONGO_ADDR")
	updateFromEnv(&MongoDatabase, "MONGO_DATABASE")
}
