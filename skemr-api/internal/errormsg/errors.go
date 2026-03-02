package errormsg

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type ErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
	Status  int               `json:"status"`
}

func (e *ErrorResponse) Error() string {
	return e.Message
}

// WriteErrorResponse is a helper function to write an error response in a consistent format.
func WriteErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	var errorResponse *ErrorResponse

	if errors.As(err, &errorResponse) {
		render.Status(r, errorResponse.Status)
		render.JSON(w, r, errorResponse)
		return
	}
	slog.Warn("Unhandled error type, returning generic error response", "error", err)

	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, ErrorResponse{
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
)
