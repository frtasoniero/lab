package controllers

import (
	"net/http"

	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/dtos"
	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/repositories"
	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/services"
	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(server *gin.Engine, repo *repositories.UserRepository) {
	service := services.NewUserService(repo)
	controller := &UserController{service: service}
	routes := server.Group("/users")
	{
		routes.POST("", controller.CreateUser)
		routes.POST("/login", controller.UserLogin)
	}
}

// @Tags users
// @Router /users [post]
// @Summary Create a new user
// @Description Register a new user on database
// @Accept json
// @Produce json
// @Param user body dtos.User true "User data"
// @Success 201 {object} dtos.Message "User created"
// @Failure 400 {object} dtos.APIError "Validation error"
// @Failure 409 {object} dtos.APIError "User already exists"
func (c *UserController) CreateUser(ginContext *gin.Context) {
	var userDto dtos.User

	err := utils.ValidateRequestBody(ginContext, &userDto)
	if err != nil {
		ginContext.Error(err)
		return
	}

	err = c.service.ExistsUserByEmail(dtos.ExistsFilter{
		Field: "email",
		Value: userDto.Email,
	})
	if err != nil {
		ginContext.Error(err)
		return
	}

	err = c.service.CreateUser(userDto.Email, userDto.Password)
	if err != nil {
		ginContext.Error(err)
		return
	}

	ginContext.JSON(http.StatusCreated, dtos.Message{
		Message: "User successfully created.",
	})
}

// @Tags users
// @Router /users/login [post]
// @Summary User login
// @Description Check if user's email and password are valid
// @Accept json
// @Produce json
// @Param user body dtos.User true "User data"
// @Success 200 {object} dtos.Token "User login"
// @Failure 400 {object} dtos.APIError "Erro while trying to login"
func (c *UserController) UserLogin(ginContext *gin.Context) {
	var userDto dtos.User
	err := utils.ValidateRequestBody(ginContext, &userDto)
	if err != nil {
		ginContext.Error(err)
		return
	}
	token, err := c.service.UserLogin(userDto.Email, userDto.Password)
	if err != nil {
		ginContext.Error(err)
		return
	}
	ginContext.JSON(http.StatusOK, dtos.Token{Token: token})
}
