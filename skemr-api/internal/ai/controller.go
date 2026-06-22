package ai

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/walmaa/skemr-api/internal/controller"
	"github.com/walmaa/skemr-api/internal/errormsg"
)

type AIController struct {
	client Model
}

func NewAIController(client Model) *AIController {
	return &AIController{client: client}
}

func (h *AIController) RegisterRoutes(r chi.Router) {
	r.Post("/ai/complete", h.complete)
}

func (h *AIController) complete(w http.ResponseWriter, r *http.Request) {
	_, ok := controller.ParseUUIDParam(w, r, "projectId")
	if !ok {
		return
	}

	result, err := h.client.Complete(r.Context(), nil, nil)

	if err != nil {
		errormsg.WriteErrorResponse(w, r, err)
		return
	}
	render.JSON(w, r, result)
	
}