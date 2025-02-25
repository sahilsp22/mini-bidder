package config

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
}

func GetPGConfig() (*Config,error) {

	err := godotenv.Load("./db.env")
    if err != nil {
        return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	return &Config{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		User: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DB: os.Getenv("DB_NAME"),
	},nil
}