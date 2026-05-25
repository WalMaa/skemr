package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/internal/errormsg"
	"github.com/walmaa/skemr-common/models"
)

func ParseUUIDParam(w http.ResponseWriter, r *http.Request, name string) (uuid.UUID, bool) {
	raw := chi.URLParam(r, name)
	if raw == "" {
		errormsg.WriteErrorResponse(w, r, &models.ErrorResponse{
			Message: "missing " + name,
			Errors:  nil,
			Status:  http.StatusBadRequest,
		})
		return uuid.Nil, false
	}

	id, err := uuid.Parse(raw)
	if err != nil {
		errormsg.WriteErrorResponse(w, r, &models.ErrorResponse{
			Message: "invalid " + name,
			Errors:  nil,
			Status:  http.StatusBadRequest,
		})
		return uuid.Nil, false
	}

	return id, true
}
