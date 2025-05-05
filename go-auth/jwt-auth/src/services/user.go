package services

import (
	"time"

	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/dtos"
	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/entities"
	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/repositories"
	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/utils"
	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/utils/middlewares"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (service *UserService) CreateUser(email, password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	user := entities.User{
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	contextServer := utils.CreateContextServerWithTimeout()
	err = service.repo.Create(contextServer, user)
	if err != nil {
		return err
	}
	return nil
}

func (service *UserService) UserLogin(bodyEmail, bodyPassword string) (string, error) {
	var user entities.User

	contextServer := utils.CreateContextServerWithTimeout()
	params := dtos.GetAnyFilter{
		Field:  "email",
		Value:  bodyEmail,
		Result: &user,
	}

	err := service.repo.GetByAny(contextServer, params)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(bodyPassword))
	if err != nil {
		return "", utils.BadRequestError("credenciais inválidas")
	}
	return middlewares.GenerateToken(user.Email, user.Id)
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", utils.BadRequestError("Error while hashing password")
	}
	return string(hashedPassword), nil
}

func (service *UserService) ExistsUserByEmail(param dtos.ExistsFilter) error {
	contextServer := utils.CreateContextServerWithTimeout()
	err := service.repo.ExistsByAny(contextServer, param)
	if err != nil {
		return err
	}
	return nil
}
