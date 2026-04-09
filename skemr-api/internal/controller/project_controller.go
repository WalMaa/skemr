package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/internal/dto"
	"github.com/walmaa/skemr-api/internal/errormsg"
	"github.com/walmaa/skemr-api/internal/service"
	"github.com/walmaa/skemr-api/internal/validation"
	"github.com/walmaa/skemr-common/models"
)

type ProjectController struct {
	Service *service.ProjectService
}

func NewProjectController(s *service.ProjectService) *ProjectController {
	return &ProjectController{Service: s}
}

func (h *ProjectController) GetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.Service.GetProjects(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(projects); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ProjectController) CreateProject(w http.ResponseWriter, r *http.Request) {
	var body dto.ProjectCreationDto

	if err := render.Decode(r, &body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := validation.Validate.Struct(body)
	if err != nil {
		errorResponse := validation.CreateErrorResponse(err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, errorResponse)
		return
	}

	user, err := h.Service.CreateProject(r.Context(), body)
	if err != nil {
		errormsg.WriteErrorResponse(w, r, err)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user)
}

func (h *ProjectController) GetProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectId"))
	if err != nil {
		errormsg.WriteErrorResponse(w, r, &models.ErrorResponse{
			Message: errormsg.ErrInvalidIdFormat,
			Status:  http.StatusBadRequest,
		})
		return
	}
	project, err := h.Service.GetProject(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, project)
}

func (h *ProjectController) DeleteProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectId"))
	if err != nil {
		errormsg.WriteErrorResponse(w, r, &models.ErrorResponse{
			Message: errormsg.ErrInvalidIdFormat,
			Status:  http.StatusBadRequest,
		})
		return
	}
	err = h.Service.DeleteProject(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
