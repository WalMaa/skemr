package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/internal/errormsg"
	"github.com/walmaa/skemr-api/internal/service"
)

type ProjectController struct {
	Service *service.ProjectService
}

func NewProjectController(s *service.ProjectService) *ProjectController {
	return &ProjectController{Service: s}
}

func (h *ProjectController) RegisterRoutes(r chi.Router) {
	r.Route("/projects", func(r chi.Router) {
		r.Post("/", h.createProject)
		r.Get("/", h.getProjects)
		r.Get("/{projectId}", h.getProject)
	})
}

func (h *ProjectController) getProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.Service.GetProjects(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func (h *ProjectController) createProject(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := h.Service.CreateProject(r.Context(), body.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *ProjectController) getProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectId"))
	if err != nil {
		http.Error(w, errormsg.ErrInvalidIdFormat.Error(), http.StatusBadRequest)
		return
	}
	project, err := h.Service.GetProject(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}
