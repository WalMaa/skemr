package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/internal/dto"
	"github.com/walmaa/skemr-api/internal/service"
)

type RuleController struct {
	Service *service.RuleService
}

func NewRuleController(s *service.RuleService) *RuleController {
	return &RuleController{Service: s}
}

func (h *RuleController) RegisterRoutes(r chi.Router) {
	r.Route("/databases/{databaseId}/rules", func(r chi.Router) {
		r.Post("/", h.createRule)
		r.Get("/", h.ListRules)
		r.Get("/{ruleId}", h.GetRule)
		r.Delete("/{ruleId}", h.deleteRule)
	})
}

func (h *RuleController) GetRule(w http.ResponseWriter, r *http.Request) {
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
	ruleId, err := uuid.Parse(chi.URLParam(r, "ruleId"))
	if err != nil {
		http.Error(w, "invalid ruleId", http.StatusBadRequest)
		return
	}
	rule, err := h.Service.GetRule(r.Context(), projectID, databaseId, ruleId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, rule)
}

func (h *RuleController) deleteRule(w http.ResponseWriter, r *http.Request) {
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
	ruleId, err := uuid.Parse(chi.URLParam(r, "ruleId"))
	if err != nil {
		http.Error(w, "invalid ruleId", http.StatusBadRequest)
		return
	}
	err = h.Service.DeleteRule(r.Context(), projectID, databaseId, ruleId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.Status(r, http.StatusNoContent)
}

func (h *RuleController) createRule(w http.ResponseWriter, r *http.Request) {
	databaseId, err := uuid.Parse(chi.URLParam(r, "databaseId"))
	if err != nil {
		http.Error(w, "invalid databaseId", http.StatusBadRequest)
		return
	}
	projectID, ok := r.Context().Value("projectId").(uuid.UUID)
	if !ok {
		http.Error(w, "projectId not found in context", http.StatusBadRequest)
		return
	}
	var body dto.RuleCreationDto
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rule, err := h.Service.CreateRule(r.Context(), projectID, databaseId, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, rule)
	render.Status(r, http.StatusCreated)
}

func (h *RuleController) ListRules(w http.ResponseWriter, r *http.Request) {
	databaseId, err := uuid.Parse(chi.URLParam(r, "databaseId"))
	if err != nil {
		http.Error(w, "invalid databaseId", http.StatusBadRequest)
		return
	}
	projectID, ok := r.Context().Value("projectId").(uuid.UUID)
	if !ok {
		http.Error(w, "projectId not found in context", http.StatusBadRequest)
		return
	}
	rules, err := h.Service.ListRulesByDatabase(r.Context(), projectID, databaseId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, rules)
}
