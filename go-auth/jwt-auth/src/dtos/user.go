package dtos

type User struct {
	Email    string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" binding:"required" example:"12345a!"`
}
