package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/walmaa/skemr-common/models"
)

var Validate *validator.Validate

func Init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func CreateErrorResponse(err error) models.ErrorResponse {
	validationErrors := err.(validator.ValidationErrors)
	errorResponse := models.ErrorResponse{
		Message: "Validation failed",
		Errors:  make(map[string]string),
	}

	for _, fieldErr := range validationErrors {
		errorResponse.Errors[strings.ToLower(fieldErr.Field())] = fieldErr.Error()
	}

	return errorResponse
}
