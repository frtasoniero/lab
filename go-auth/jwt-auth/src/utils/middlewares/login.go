package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var SecretKey = []byte("string-secret-to-password-at-least-256-bits-long")

type Claims struct {
	Email  string `json:"email"`
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

func GenerateToken(email string, userId primitive.ObjectID) (string, error) {
	expirationTime := utils.TimeNowBrazil().Add(24 * time.Hour)

	claims := &Claims{
		Email:  email,
		UserId: userId.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", utils.BadRequestError("Error while generating token")
	}

	return tokenString, nil
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" && (c.Request.URL.Path == "/users" || c.Request.URL.Path == "/users/login") ||
			c.Request.Method == "GET" && strings.HasPrefix(c.Request.URL.Path, "/swagger/") {
			c.Next()
			return
		}

		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found or Token is not correct"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}

		objUserId, err := primitive.ObjectIDFromHex(claims.UserId)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid userId"})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("userId", objUserId)
		c.Next()
	}
}
