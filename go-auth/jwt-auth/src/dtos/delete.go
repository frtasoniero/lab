package dtos

import "go.mongodb.org/mongo-driver/bson/primitive"

type DeleteFilter struct {
	Id              primitive.ObjectID
	ForeignKey      string
	ForeignKeyValue primitive.ObjectID
}
