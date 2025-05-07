package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	MongoURI  string
	MongoDB   string
	SecretKey []byte
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	MongoURI = getEnv("MONGO_URI", "mongodb://admin:secret@localhost:27017")
	MongoDB = getEnv("MONGO_DB", "default-db")
	SecretKey = []byte(getEnv("SECRET_KEY", "string-secret-to-password-at-least-256-bits-long"))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
