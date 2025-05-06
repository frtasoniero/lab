package dtos

import "go.mongodb.org/mongo-driver/bson/primitive"

type UptadeFilter struct {
	Id              primitive.ObjectID
	Dto             interface{}
	ForeignKey      string
	ForeignKeyValue primitive.ObjectID
}
