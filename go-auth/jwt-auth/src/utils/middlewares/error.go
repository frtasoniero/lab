package middlewares

import (
	"net/http"

	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/dtos"
	"github.com/gin-gonic/gin"
)

func ErrorMiddlewareHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			apiErr, ok := err.(*dtos.APIError)
			if ok {
				c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
				c.Abort()
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server internal error"})
			c.Abort()
		}
	}
}
