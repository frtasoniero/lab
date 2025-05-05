package converts

import (
	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func StringToObject(value string) (primitive.ObjectID, error) {
	valueObj, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		return primitive.NilObjectID, utils.BadRequestError("Error while trying to convert string to object")
	}
	return valueObj, err
}
