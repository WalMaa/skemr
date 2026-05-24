package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/internal/errormsg"
	"github.com/walmaa/skemr-api/internal/service"
)

type DatabaseChangeController struct {
	Service *service.DatabaseChangeService
}

func NewDatabaseChangeController(s *service.DatabaseChangeService) *DatabaseChangeController {
	return &DatabaseChangeController{Service: s}
}

func (h *DatabaseChangeController) RegisterRoutes(r chi.Router) {
	r.Route("/databases/{databaseId}/changes", func(r chi.Router) {

		r.Get("/", h.listDatabaseChanges)
		r.Get("/{databaseChangeId}", h.GetDatabaseChange)
	})
}

func (h *DatabaseChangeController) listDatabaseChanges(w http.ResponseWriter, r *http.Request) {
	projectID, ok := r.Context().Value("projectId").(uuid.UUID)
	if !ok {
		http.Error(w, "projectId not found in context", http.StatusBadRequest)
		return
	}
	databaseId, err := uuid.Parse(chi.URLParam(r, "databaseId"))
	if err != nil {
		http.Error(w, "invalid databaseId", http.StatusBadRequest)
		return
	}

	databaseChanges, err := h.Service.GetDatabaseChanges(r.Context(), projectID, databaseId, 100, 0)

	if err != nil {
		errormsg.WriteErrorResponse(w, r, err)
		return
	}
	render.JSON(w, r, databaseChanges)
}

func (h *DatabaseChangeController) GetDatabaseChange(w http.ResponseWriter, r *http.Request) {
	projectID, ok := r.Context().Value("projectID").(uuid.UUID)
	if !ok {
		http.Error(w, "projectId not found in context", http.StatusBadRequest)
		return
	}
	databaseId, err := uuid.Parse(chi.URLParam(r, "databaseId"))
	if err != nil {
		http.Error(w, "invalid databaseId", http.StatusBadRequest)
		return
	}

	databaseChangeId, err := uuid.Parse(chi.URLParam(r, "databaseChangeId"))
	if err != nil {
		http.Error(w, "invalid databaseChangeId", http.StatusBadRequest)
	}

	databaseChange, err := h.Service.GetDatabaseChange(r.Context(), projectID, databaseId, databaseChangeId)

	if err != nil {
		errormsg.WriteErrorResponse(w, r, err)
		return
	}
	render.JSON(w, r, databaseChange)
}
