package services

import (
	"time"

	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/dtos"
	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/entities"
	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/repositories"
	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskService struct {
	repo *repositories.TaskRepository
}

func NewTaskService(repo *repositories.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (service *TaskService) CreateTask(userId primitive.ObjectID, title, description string) error {
	task := entities.Task{
		Title:       title,
		Description: description,
		UserId:      userId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	contextServer := utils.CreateContextServerWithTimeout()
	err := service.repo.Create(contextServer, task)
	if err != nil {
		return err
	}

	return nil
}

func (service *TaskService) ExistsTaskByTitle(param dtos.ExistsFilter) error {
	contextServer := utils.CreateContextServerWithTimeout()
	err := service.repo.ExistsByAny(contextServer, param)
	if err != nil {
		return err
	}
	return nil
}

func (service *TaskService) GetTaskById(params dtos.GetAnyFilter) error {
	contextServer := utils.CreateContextServerWithTimeout()
	return service.repo.GetByAny(contextServer, params)
}
