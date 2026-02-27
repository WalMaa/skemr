package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func Init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

type ErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

func CreateErrorResponse(err error) ErrorResponse {
	validationErrors := err.(validator.ValidationErrors)
	errorResponse := ErrorResponse{
		Message: "Validation failed",
		Errors:  make(map[string]string),
	}

	for _, fieldErr := range validationErrors {
		errorResponse.Errors[strings.ToLower(fieldErr.Field())] = fieldErr.Error()
	}

	return errorResponse
}
