package utils

import (
	"net/http"

	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/dtos"
)

func NotFoundError(message string) *dtos.APIError {
	return &dtos.APIError{StatusCode: http.StatusNotFound, Message: message}
}
func BadRequestError(message string) *dtos.APIError {
	return &dtos.APIError{StatusCode: http.StatusBadRequest, Message: message}
}
func ConflictError(message string) *dtos.APIError {
	return &dtos.APIError{StatusCode: http.StatusConflict, Message: message}
}
func InternalServerError(message string) *dtos.APIError {
	return &dtos.APIError{StatusCode: http.StatusInternalServerError, Message: message}
}
