package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/internal/errormsg"
	"github.com/walmaa/skemr-api/internal/service"
	"github.com/walmaa/skemr-common/models"
)

type IntegrationController struct {
	IntegrationService *service.IntegrationService
}

func NewIntegrationController(s *service.IntegrationService) *IntegrationController {
	return &IntegrationController{IntegrationService: s}
}

func (h *IntegrationController) RegisterRoutes(r chi.Router) {
	r.Get("/ci-cd/rules", h.listRulesByDatabase)

}

func (h *IntegrationController) listRulesByDatabase(w http.ResponseWriter, r *http.Request) {
	databaseId, err := uuid.Parse(chi.URLParam(r, "databaseId"))
	if err != nil {
		errormsg.WriteErrorResponse(w, r, &models.ErrorResponse{
			Message: "Invalid database ID format",
			Status:  http.StatusBadRequest,
			Errors:  nil,
		},
		)
		return
	}

	projectId, err := uuid.Parse(chi.URLParam(r, "projectId"))
	if err != nil {
		errormsg.WriteErrorResponse(w, r, &models.ErrorResponse{
			Message: "Invalid project ID format",
			Status:  http.StatusBadRequest,
			Errors:  nil,
		},
		)
		return
	}

	rules, err := h.IntegrationService.ListRulesByDatabase(r.Context(), projectId, databaseId)
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
