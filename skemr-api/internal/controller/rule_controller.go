package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/walmaa/skemr-api/internal/dto"
	"github.com/walmaa/skemr-api/internal/errormsg"
	"github.com/walmaa/skemr-api/internal/service"
	"github.com/walmaa/skemr-api/internal/validation"
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
	projectId, ok := ParseUUIDParam(w, r, "projectId")
	if !ok {
		return
	}

	databaseId, ok := ParseUUIDParam(w, r, "databaseId")
	if !ok {
		return
	}

	ruleId, ok := ParseUUIDParam(w, r, "ruleId")
	if !ok {
		return
	}

	rule, err := h.Service.GetRule(r.Context(), projectId, databaseId, ruleId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, rule)
}

func (h *RuleController) deleteRule(w http.ResponseWriter, r *http.Request) {
	projectId, ok := ParseUUIDParam(w, r, "projectId")
	if !ok {
		return
	}

	databaseId, ok := ParseUUIDParam(w, r, "databaseId")
	if !ok {
		return
	}

	ruleId, ok := ParseUUIDParam(w, r, "ruleId")
	if !ok {
		return
	}

	err := h.Service.DeleteRule(r.Context(), projectId, databaseId, ruleId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.Status(r, http.StatusNoContent)
}

func (h *RuleController) createRule(w http.ResponseWriter, r *http.Request) {
	projectId, ok := ParseUUIDParam(w, r, "projectId")
	if !ok {
		return
	}

	databaseId, ok := ParseUUIDParam(w, r, "databaseId")
	if !ok {
		return
	}

	var body dto.RuleCreationDto

	if err := render.Decode(r, &body); err != nil {
		errormsg.WriteErrorResponse(w, r, err)
		return
	}

	err := validation.Validate.Struct(body)

	if err != nil {
		errorResponse := validation.CreateErrorResponse(err)
		errormsg.WriteErrorResponse(w, r, &errorResponse)
		return
	}

	rule, err := h.Service.CreateRule(r.Context(), projectId, databaseId, body)

	if err != nil {
		errormsg.WriteErrorResponse(w, r, err)
		return
	}

	render.JSON(w, r, rule)
	render.Status(r, http.StatusCreated)
}

func (h *RuleController) ListRules(w http.ResponseWriter, r *http.Request) {
	projectId, ok := ParseUUIDParam(w, r, "projectId")
	if !ok {
		return
	}

	databaseId, ok := ParseUUIDParam(w, r, "databaseId")
	if !ok {
		return
	}

	rules, err := h.Service.ListRulesByDatabase(r.Context(), projectId, databaseId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, rules)
}
