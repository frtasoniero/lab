package controllers

import (
	"net/http"

	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/dtos"
	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/entities"
	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/repositories"
	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/services"
	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/utils"
	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/utils/converts"
	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/utils/enums"
	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/utils/formats"
	"github.com/gin-gonic/gin"
)

type TaskController struct {
	service *services.TaskService
}

func NewTaskController(server *gin.Engine, repo *repositories.TaskRepository) {
	service := services.NewTaskService(repo)
	controller := &TaskController{
		service: service,
	}
	routes := server.Group("/tasks")
	{
		routes.POST("", controller.CreateTask)
		routes.GET("/:id", controller.GetTaskById)
	}
}

// @Security BearerAuth
// @Tags tasks
// @Router /tasks [post]
// @Summary Create new task
// @Description Register a new task on database
// @Accept json
// @Produce json
// @Param task body dtos.Task true "Task data"
// @Success 201 {object} dtos.Message "Task created"
// @Failure 400 {object} dtos.APIError "Validation error"
// @Failure 409 {object} dtos.APIError "Task already exists"
func (c *TaskController) CreateTask(ginContext *gin.Context) {
	var taskDto dtos.Task

	err := utils.ValidateRequestBody(ginContext, &taskDto)
	if err != nil {
		ginContext.Error(err)
		return
	}

	userId, err := utils.GetUserAuthenticated(ginContext)
	if err != nil {
		ginContext.Error(err)
		return
	}

	err = c.service.ExistsTaskByTitle(dtos.ExistsFilter{
		Field:           "title",
		Value:           taskDto.Title,
		ForeignKey:      "userId",
		ForeignKeyValue: userId,
	})
	if err != nil {
		ginContext.Error(err)
		return
	}

	err = c.service.CreateTask(userId, taskDto.Title, taskDto.Description)
	if err != nil {
		ginContext.Error(err)
		return
	}

	ginContext.JSON(http.StatusOK, dtos.Message{
		Message: "Task successfully created.",
	})
}

// @Security BearerAuth
// @Tags tasks
// @Router /tasks/{id} [get]
// @Summary Get task by id
// @Description Get a task by its id from database
// @Accept json
// @Produce json
// @Param id path string true "Task Id" example("6817d32b97e0903fae78beec")
// @Success 200 {object} dtos.Message "Task by id"
// @Failure 400 {object} dtos.APIError "Validation error"
func (c *TaskController) GetTaskById(ginContext *gin.Context) {
	var task entities.Task
	idHex := ginContext.Param("id")

	id, err := converts.StringToObject(idHex)
	if err != nil {
		ginContext.Error(err)
		return
	}

	userId, err := utils.GetUserAuthenticated(ginContext)
	if err != nil {
		ginContext.Error(err)
		return
	}

	params := dtos.GetAnyFilter{
		Field:           "_id",
		Value:           id,
		ForeignKey:      "userId",
		ForeignKeyValue: userId,
		Result:          &task,
	}

	err = c.service.GetTaskById(params)
	if err != nil {
		ginContext.Error(err)
		return
	}

	createdAt := utils.TimeBrazil(task.CreatedAt)
	updatedAt := utils.TimeBrazil(task.UpdatedAt)

	task.CreatedAt, err = formats.Time(createdAt, enums.FormatTime.DataHour())
	if err != nil {
		ginContext.Error(err)
		return
	}

	task.UpdatedAt, err = formats.Time(updatedAt, enums.FormatTime.DataHour())
	if err != nil {
		ginContext.Error(err)
		return
	}

	ginContext.JSON(http.StatusOK, task)
}
