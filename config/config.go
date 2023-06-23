package config

import (
	"log"

	"github.com/joho/godotenv"
)

type DbConfig struct {
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`
	ServerPort     string `mapstructure:"PORT"`

	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`
}

func LoadConfig(path string) (DbConfig, error) {

	env, e := godotenv.Read(path)
	if e != nil {
		log.Fatalf("Error loading .env file")
	}

	conf := DbConfig{
		DBHost:         env["POSTGRES_HOST"],
		DBUserName:     env["POSTGRES_USER"],
		DBUserPassword: env["POSTGRES_PASSWORD"],
		DBName:         env["POSTGRES_DB"],
		DBPort:         env["POSTGRES_PORT"],
		ServerPort:     env["POSTGRES_PORT"],
		ClientOrigin:   env["CLIENT_ORIGIN"],
	}

	return conf, nil
}
