package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/frtasoniero/lab/go-auth/pwh-auth/internal/auth"
	"github.com/frtasoniero/lab/go-auth/pwh-auth/internal/db"
)

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // fallback
	}

	if err := db.Connect(mongoURI); err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	users := []auth.User{
		{Username: "alice", PasswordHash: hash("alice123")},
		{Username: "bob", PasswordHash: hash("bobsecure")},
		{Username: "test", PasswordHash: hash("testpass")},
	}

	collection := db.MongoClient.Database("pwh_db").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, user := range users {
		_, err := collection.InsertOne(ctx, user)
		if err != nil {
			log.Printf("Error inserting %s: %v", user.Username, err)
		} else {
			log.Printf("Inserted test user: %s", user.Username)
		}
	}
}

func hash(password string) string {
	hashed, err := auth.HashPassword(password)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}
	return hashed
}
