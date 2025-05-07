package main

import (
	"log"

	_ "github.com/frtasoniero/lab/go-auth/jwt-auth/docs"

	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/config"
	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/controllers"
	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/repositories"
	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/utils/middlewares"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	uri    = "mongodb://admin:secret123@localhost:27017"
	dbName = "jwt-auth-local"
)

// @title JWT Auth API
// @version 1.0
// @description Authentication API with JWT
// @host localhost:8080
// @BasePath /
func main() {
	config.LoadEnv()

	log.Println("DB:", config.MongoDB)
	log.Println("URI:", config.MongoURI)

	repoUser, errUser := repositories.NewUserRepository(uri, dbName, "users")
	repoTask, errTask := repositories.NewTaskRepository(uri, dbName, "tasks")

	if errUser != nil || errTask != nil {
		log.Fatalf("Error while instantiating repositories: errUser=%v, errTask=%v", errUser, errTask)
		return
	}

	server := gin.Default()
	server.Use(middlewares.CorsMiddleware())
	server.Use(middlewares.ErrorMiddlewareHandler())
	server.Use(middlewares.JWTAuthMiddleware())

	controllers.NewUserController(server, repoUser)
	controllers.NewTaskController(server, repoTask)

	// @securityDefinitions.apikey BearerAuth
	// @in header
	// @name Authorization
	// @description Value: Bearer abc... (Bearer+space+token)
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.DefaultModelsExpandDepth(-1)))

	// Initialize server on port 8080
	server.Run(":8080")
}
