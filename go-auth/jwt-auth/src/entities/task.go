package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId      primitive.ObjectID `bson:"userId" json:"userId"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
