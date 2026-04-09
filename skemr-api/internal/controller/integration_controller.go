package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/internal/service"
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
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	projectId, err := uuid.Parse(chi.URLParam(r, "projectId"))
	if err != nil {
		http.Error(w, "Invalid project ID format", http.StatusBadRequest)
		return
	}

	rules, err := h.IntegrationService.ListRulesByDatabase(r.Context(), projectId, databaseId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, rules)
}
