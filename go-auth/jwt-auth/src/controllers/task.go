package controllers

import (
	"net/http"

	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/dtos"
	dtosPage "github.com/frtasoniero/lab/go-auth/jwt-auth/src/dtos/pagination"
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
		routes.GET("/", controller.GetTasks)
		routes.PUT("/:id", controller.UptadeTask)
		routes.DELETE("/:id", controller.DeleteTask)
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

// @Security BearerAuth
// @Tags tasks
// @Router /tasks [get]
// @Summary Get task pagination
// @Description Get task pagination from the API
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limitPage query int false "Number of records per page"
// @Param searchField query string false "Search by which task property"
// @Param searchValue query string false "Search by value of the task property"
// @Param sortField query string false "Sort by which task property"
// @Param sortOrder query string false "Sorting" Enums(ascending, descending)
// @Success 200 {array} dtosPage.PaginationResult[entities.Task] "Task list"
// @Failure 400 {object} dtos.APIError "Validation error"
func (c *TaskController) GetTasks(ginContext *gin.Context) {
	page := ginContext.DefaultQuery("page", "1")
	limitPage := ginContext.DefaultQuery("limitPage", "5")
	searchField := ginContext.DefaultQuery("searchField", "")
	searchValue := ginContext.DefaultQuery("searchValue", "")
	sortField := ginContext.DefaultQuery("sortField", "_id")
	sortOrderStr := ginContext.DefaultQuery("sortOrder", enums.SortOrder.AscendingStr())
	sortOrder, err := enums.SortOrder.ConvertSortOrderEnumToInt(sortOrderStr)
	if err != nil {
		ginContext.Error(err)
		return
	}
	skipInt, err := converts.StringToInt(page)
	if err != nil {
		ginContext.JSON(http.StatusBadRequest, dtos.APIError{Message: "Invalid skip parameter"})
		return
	}
	limitInt, err := converts.StringToInt(limitPage)
	if err != nil {
		ginContext.JSON(http.StatusBadRequest, dtos.APIError{Message: "Invalid perPage parameter"})
		return
	}
	userId, err := utils.GetUserAuthenticated(ginContext)
	if err != nil {
		ginContext.Error(err)
		return
	}
	var tasks []entities.Task
	paginationResult := c.service.GetPagination(dtosPage.PaginationParams{
		Field:       "userId",
		Value:       userId,
		Result:      &tasks,
		Skip:        skipInt,
		Limit:       limitInt,
		SearchField: searchField,
		SearchValue: searchValue,
		SortField:   sortField,
		SortOrder:   sortOrder,
	})
	for item := range tasks {
		tasks[item].CreatedAt = utils.TimeBrazil(tasks[item].CreatedAt)
		tasks[item].UpdatedAt = utils.TimeBrazil(tasks[item].UpdatedAt)
	}
	if paginationResult.PaginationResultContext.Err != nil {
		ginContext.Error(err)
		return
	}
	ginContext.JSON(http.StatusOK, paginationResult)
}

// @Security BearerAuth
// @Tags tasks
// @Router /tasks/{id} [put]
// @Summary Edit task by Id
// @Description Edit a task by its Id
// @Accept json
// @Produce json
// @Param id path string true "Task Id" example("60c72b2f9b1d8b57b8ed2123")
// @Param task body dtos.Task true "User data"
// @Success 200 {object} dtos.Message "Task by id"
// @Failure 400 {object} dtos.APIError "Validation error"
func (c *TaskController) UptadeTask(ginContext *gin.Context) {
	var updateData dtos.Task
	err := utils.ValidateRequestBody(ginContext, &updateData)
	if err != nil {
		ginContext.Error(err)
		return
	}

	taskIdStr := ginContext.Param("id")
	taskId, err := converts.StringToObject(taskIdStr)
	if err != nil {
		ginContext.Error(err)
		return
	}

	userId, err := utils.GetUserAuthenticated(ginContext)
	if err != nil {
		ginContext.Error(err)
		return
	}

	err = c.service.UptadeTask(taskId, userId, updateData)
	if err != nil {
		ginContext.Error(err)
		return
	}
	ginContext.JSON(http.StatusOK, dtos.Message{
		Message: "Task successfully updated.",
	})
}

// @Security BearerAuth
// @Tags tasks
// @Router /tasks/{id} [delete]
// @Summary Delete task by Id
// @Description Delete a task by its Id
// @Accept json
// @Produce json
// @Param id path string true "ID da Task" example("60c72b2f9b1d8b57b8ed2123")
// @Success 200 {object} dtos.Message "Task by id"
// @Failure 400 {object} dtos.APIError "Validation error"
func (c *TaskController) DeleteTask(ginContext *gin.Context) {
	id := ginContext.Param("id")
	idObj, err := converts.StringToObject(id)
	if err != nil {
		ginContext.Error(err)
		return
	}
	userId, err := utils.GetUserAuthenticated(ginContext)
	if err != nil {
		ginContext.Error(err)
		return
	}
	params := dtos.DeleteFilter{
		Id:              idObj,
		ForeignKey:      "userId",
		ForeignKeyValue: userId,
	}
	err = c.service.DeleteTask(params)
	if err != nil {
		ginContext.Error(err)
		return
	}
	ginContext.JSON(http.StatusOK, dtos.Message{
		Message: "Task deletado com sucesso.",
	})
}
