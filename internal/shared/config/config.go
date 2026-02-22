package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ApplicationName string
	HTTPPort        string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, using system env")
	}

	return &Config{
		ApplicationName: os.Getenv("APPLICATION_NAME"),
		HTTPPort:        os.Getenv("APPLICATION_HTTPPORT"),
	}
}
