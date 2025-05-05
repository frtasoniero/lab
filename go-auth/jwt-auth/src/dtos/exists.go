package dtos

import "go.mongodb.org/mongo-driver/bson/primitive"

type ExistsFilter struct {
	Field           string
	Value           interface{}
	ForeignKey      string
	ForeignKeyValue primitive.ObjectID
}
