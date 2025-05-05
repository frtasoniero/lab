package utils

import "github.com/gin-gonic/gin"

func ValidateRequestBody(ginContext *gin.Context, dto interface{}) error {
	err := ginContext.ShouldBindJSON(&dto)
	if err != nil {
		return BadRequestError("Dados inválidos. Preencha os campos corretamente.")
	}
	return nil
}
