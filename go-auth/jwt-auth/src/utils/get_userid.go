package utils

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserAuthenticated(c *gin.Context) (primitive.ObjectID, error) {
	userId, exists := c.Get("userId")
	if !exists {
		return primitive.NilObjectID, BadRequestError("User not authenticated.")
	}
	return userId.(primitive.ObjectID), nil
}
