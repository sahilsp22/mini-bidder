package config

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
)

const (
	CACHE_TIMEOUT = 11
	CACHE_UPDATE_INTERVAL = 10
)

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
}

type Memcache struct {
	Host string
	Port string
}

func GetPGConfig() (*Postgres,error) {

	err := godotenv.Load("/Users/sahilpatil/mini-bidder/config/db.env")
    if err != nil {
        return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	return &Postgres{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		User: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DB: os.Getenv("DB_NAME"),
	},nil
}

func GetMcCConfig() (*Memcache,error) {
	err := godotenv.Load("/Users/sahilpatil/mini-bidder/config/db.env")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	return &Memcache{
		Host: os.Getenv("MC_HOST"),
		Port: os.Getenv("MC_PORT"),
	},nil
}