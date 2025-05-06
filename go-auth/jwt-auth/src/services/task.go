package services

import (
	"time"

	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/dtos"
	dtoPages "github.com/frtasoniero/lab/go-auth/jwt-auth/src/dtos/pagination"
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

func (s *TaskService) GetPagination(params dtoPages.PaginationParams) dtoPages.PaginationResult[entities.Task] {
	contextServer := utils.CreateContextServerWithTimeout()
	resultContext := s.repo.GetPagination(contextServer, params)

	return dtoPages.PaginationResult[entities.Task]{
		Items:                   *params.Result.(*[]entities.Task),
		PageCurrent:             params.Skip,
		PaginationResultContext: resultContext,
	}
}

func (service *TaskService) UptadeTask(id, userId primitive.ObjectID, dto dtos.Task) error {
	contextServer := utils.CreateContextServerWithTimeout()
	task := entities.Task{
		Title:       dto.Title,
		Description: dto.Description,
		UpdatedAt:   utils.TimeNowBrazil(),
	}
	params := dtos.UptadeFilter{
		Id:              id,
		Dto:             task,
		ForeignKey:      "userId",
		ForeignKeyValue: userId,
	}
	err := service.repo.Update(contextServer, params)
	if err != nil {
		return err
	}
	return nil
}
