package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/internal/errormsg"
	"github.com/walmaa/skemr-api/internal/service"
)

type PipelineRunController struct {
	Service *service.PipelineRunService
}

func NewPipelineRunController(s *service.PipelineRunService) *PipelineRunController {
	return &PipelineRunController{Service: s}
}

func (h *PipelineRunController) RegisterRoutes(r chi.Router) {
	r.Route("/databases/{databaseId}/pipeline-runs", func(r chi.Router) {
		r.Get("/", h.listPipelineRuns)
		r.Get("/{pipelineRunId}", h.getPipelineRun)
	})
}

func (h *PipelineRunController) listPipelineRuns(w http.ResponseWriter, r *http.Request) {
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

	pipelineRuns, err := h.Service.GetPipelineRuns(r.Context(), projectID, databaseId)

	if err != nil {
		errormsg.WriteErrorResponse(w, r, err)
		return
	}
	render.JSON(w, r, pipelineRuns)
}

func (h *PipelineRunController) getPipelineRun(w http.ResponseWriter, r *http.Request) {
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

	pipelineRunId, err := uuid.Parse(chi.URLParam(r, "pipelineRunId"))

	if err != nil {
		http.Error(w, "invalid pipelineRunId", http.StatusBadRequest)
		return
	}

	pipelineRun, err := h.Service.GetPipelineRun(r.Context(), projectID, databaseId, pipelineRunId)
	if err != nil {
		errormsg.WriteErrorResponse(w, r, err)
	}

	render.JSON(w, r, pipelineRun)

}
