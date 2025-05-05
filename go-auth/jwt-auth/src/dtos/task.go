package dtos

type Task struct {
	Title       string `json:"title" binding:"required" example:"New Task"`
	Description string `json:"description" example:"This is an example task."`
}
