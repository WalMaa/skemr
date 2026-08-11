package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/walmaa/skemr-api/internal/domain/pipelinerun"
	"github.com/walmaa/skemr-api/internal/errormsg"
	"github.com/walmaa/skemr-api/internal/requestparams"
	"github.com/walmaa/skemr-api/internal/service"
	"github.com/walmaa/skemr-api/internal/validation"
	"github.com/walmaa/skemr-common/models"
)

type IntegrationController struct {
	integrationService *service.IntegrationService
}

func NewIntegrationController(s *service.IntegrationService) *IntegrationController {
	return &IntegrationController{integrationService: s}
}

func (h *IntegrationController) RegisterRoutes(r chi.Router) {
	r.Get("/ci-cd/rules", h.listRulesByDatabase)
	r.Post("/ci-cd/pipeline-runs", h.createPipelineRun)

}

func (h *IntegrationController) listRulesByDatabase(w http.ResponseWriter, r *http.Request) {
	projectId, ok := requestparams.ParseUUIDParam(w, r, "projectId")
	if !ok {
		return
	}

	databaseId, ok := requestparams.ParseUUIDParam(w, r, "databaseId")
	if !ok {
		return
	}

	rules, err := h.integrationService.ListRulesByDatabase(r.Context(), projectId, databaseId)
	if err != nil {
		errormsg.WriteErrorResponse(w, r, &models.ErrorResponse{
			Message: "Error fetching rules",
			Status:  http.StatusInternalServerError,
			Errors:  nil,
		},
		)
		return
	}
	render.JSON(w, r, rules)
}

func (h *IntegrationController) createPipelineRun(w http.ResponseWriter, r *http.Request) {

	projectId, ok := requestparams.ParseUUIDParam(w, r, "projectId")
	if !ok {
		return
	}

	databaseId, ok := requestparams.ParseUUIDParam(w, r, "databaseId")
	if !ok {
		return
	}

	var req pipelinerun.PipelineRunCreationDto
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		errormsg.WriteInvalidRequestBodyErrorResponse(w, r)
		return
	}

	err := validation.Validate.Struct(req)

	if err != nil {
		errorResponse := validation.CreateErrorResponse(err)
		errormsg.WriteErrorResponse(w, r, &errorResponse)
		return
	}

	rule, err := h.integrationService.CreatePipeLineRun(r.Context(), projectId, databaseId, req)
	if err != nil {
		errormsg.WriteErrorResponse(w, r, err)
		return
	}
	render.JSON(w, r, rule)
}
