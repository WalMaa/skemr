package errormsg

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/walmaa/skemr-common/models"
)

// WriteErrorResponse is a helper function to write an error response in a consistent format.
func WriteErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	var errorResponse *models.ErrorResponse

	if errors.As(err, &errorResponse) {
		render.Status(r, errorResponse.Status)
		render.JSON(w, r, errorResponse)
		return
	}
	slog.Warn("Unhandled error type, returning generic error response", "error", err)

	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, models.ErrorResponse{
		Message: "Internal Server Error",
		Status:  http.StatusInternalServerError,
	},
	)
}

var (
	ErrDatabaseAlreadyExists = "database already exists"
	ErrDatabaseNotFound      = "database not found"
	ErrProjectNotFound       = "project not found"
	ErrInvalidIdFormat       = "invalid id format"
	ErrExpiryTimeInPast      = "expiry time is in the past"
	ErrRuleWithSameName      = "rule with the same name already exists"
)
