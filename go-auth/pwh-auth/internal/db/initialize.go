// internal/db/initialize.go
package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitSchemas creates collections and indexes (shared logic)
func InitSchemas(databaseName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := MongoClient.Database(databaseName)
	users := db.Collection("users")

	// Unique index on "username"
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("unique_username"),
	}

	_, err := users.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("Failed to create index in DB %s: %v", databaseName, err)
		return err
	}

	log.Printf("Database schema initialized in %s", databaseName)
	return nil
}
