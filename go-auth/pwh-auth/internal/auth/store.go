package auth

import (
	"context"
	"time"

	"github.com/frtasoniero/lab/go-auth/pwt-auth/internal/db"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	Username     string `bson:"username"`
	PasswordHash string `bson:"password_hash"`
}

func SaveUser(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	colletion := db.MongoClient.Database("pwh_db").Collection("users")
	_, err := colletion.InsertOne(ctx, user)
	return err
}

func GetUserByUsername(username string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	collection := db.MongoClient.Database("pwh_db").Collection("users")
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
