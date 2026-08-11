package databasechange

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/walmaa/skemr-api/internal/errormsg"
	"github.com/walmaa/skemr-api/internal/requestparams"
)

type DatabaseChangeController struct {
	Service *DatabaseChangeService
}

func NewDatabaseChangeController(s *DatabaseChangeService) *DatabaseChangeController {
	return &DatabaseChangeController{Service: s}
}

func (h *DatabaseChangeController) RegisterRoutes(r chi.Router) {
	r.Route("/databases/{databaseId}/changes", func(r chi.Router) {

		r.Get("/", h.listDatabaseChanges)
		r.Get("/{databaseChangeId}", h.GetDatabaseChange)
	})
}

func (h *DatabaseChangeController) listDatabaseChanges(w http.ResponseWriter, r *http.Request) {
	projectId, ok := requestparams.ParseUUIDParam(w, r, "projectId")
	if !ok {
		return
	}

	databaseId, ok := requestparams.ParseUUIDParam(w, r, "databaseId")
	if !ok {
		return
	}

	databaseChanges, err := h.Service.GetDatabaseChanges(r.Context(), projectId, databaseId, 100, 0)

	if err != nil {
		errormsg.WriteErrorResponse(w, r, err)
		return
	}
	render.JSON(w, r, databaseChanges)
}

func (h *DatabaseChangeController) GetDatabaseChange(w http.ResponseWriter, r *http.Request) {
	projectId, ok := requestparams.ParseUUIDParam(w, r, "projectId")
	if !ok {
		return
	}

	databaseId, ok := requestparams.ParseUUIDParam(w, r, "databaseId")
	if !ok {
		return
	}

	databaseChangeId, ok := requestparams.ParseUUIDParam(w, r, "databaseChangeId")
	if !ok {
		return
	}

	databaseChange, err := h.Service.GetDatabaseChange(r.Context(), projectId, databaseId, databaseChangeId)

	if err != nil {
		errormsg.WriteErrorResponse(w, r, err)
		return
	}
	render.JSON(w, r, databaseChange)
}
