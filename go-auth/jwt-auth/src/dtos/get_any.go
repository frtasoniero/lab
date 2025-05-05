package dtos

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetAnyFilter struct {
	Field           string
	Value           interface{}
	Result          interface{}
	ForeignKey      string
	ForeignKeyValue primitive.ObjectID
}
