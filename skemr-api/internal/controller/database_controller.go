package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/internal/dto"
	"github.com/walmaa/skemr-api/internal/service"
)

type DatabaseController struct {
	Service *service.DatabaseService
}

func NewDatabaseController(s *service.DatabaseService) *DatabaseController {
	return &DatabaseController{Service: s}
}

func (h *DatabaseController) RegisterRoutes(r chi.Router) {
	r.Route("/databases", func(r chi.Router) {
		r.Post("/", h.createDatabase)
		r.Get("/", h.listDatabasesByProject)
		r.Get("/{databaseId}", h.getDatabase)
		r.Delete("/{databaseId}", h.deleteDatabase)
		r.Patch("/{databaseId}", h.updateDatabase)
		r.Post("/{databaseId}/sync", h.syncDatabase)
	})
}

func (h *DatabaseController) deleteDatabase(w http.ResponseWriter, r *http.Request) {
	databaseId, err := uuid.Parse(chi.URLParam(r, "databaseId"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	err = h.Service.DeleteDatabase(r.Context(), databaseId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *DatabaseController) listDatabasesByProject(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "projectId"))
	if err != nil {
		http.Error(w, "Invalid project ID format", http.StatusBadRequest)
		return
	}
	databases, err := h.Service.ListDatabasesByProject(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(databases); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *DatabaseController) createDatabase(w http.ResponseWriter, r *http.Request) {
	projectId, err := uuid.Parse(chi.URLParam(r, "projectId"))
	if err != nil {
		http.Error(w, "Invalid project ID format", http.StatusBadRequest)
		return
	}
	var body dto.DatabaseCreationDto

	err = render.Decode(r, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validate.Struct(body)

	database, err := h.Service.CreateDatabase(r.Context(), projectId, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, database)
}

func (h *DatabaseController) updateDatabase(w http.ResponseWriter, r *http.Request) {
	projectId := r.Context().Value("projectId").(uuid.UUID)
	databaseId, err := uuid.Parse(chi.URLParam(r, "databaseId"))
	if err != nil {
		http.Error(w, "Invalid database ID format", http.StatusBadRequest)
		return
	}
	var body dto.DatabaseUpdateDto
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	database, err := h.Service.UpdateDatabase(r.Context(), projectId, databaseId, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(database); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *DatabaseController) syncDatabase(w http.ResponseWriter, r *http.Request) {
	projectId := r.Context().Value("projectId").(uuid.UUID)
	databaseId, err := uuid.Parse(chi.URLParam(r, "databaseId"))
	if err != nil {
		http.Error(w, "Invalid database ID format", http.StatusBadRequest)
		return
	}
	err = h.Service.EnqueueManualDatabaseSync(r.Context(), projectId, databaseId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *DatabaseController) getDatabase(w http.ResponseWriter, r *http.Request) {
	databaseId, err := uuid.Parse(chi.URLParam(r, "databaseId"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	database, err := h.Service.GetDatabase(r.Context(), databaseId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(database); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
