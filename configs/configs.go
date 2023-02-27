package configs

import (
	"os"
)

type DbConfig struct {
	Hosts     string
	DbName    string
	DbOptions string
}

var MongoConfig DbConfig

func SetupDbConfigs() *DbConfig {
	MongoConfig = DbConfig{
		Hosts:     os.Getenv("MONGO_ENV"),
		DbName:    os.Getenv("MONGO_DATABSE"),
		DbOptions: os.Getenv("MONGO_OPTIONS"),
	}
	return &MongoConfig
}
