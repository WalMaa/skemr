package ai

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
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
	projectId, ok := controller.ParseUUIDParam(w, r, "projectId")
	if !ok {
		return
	}

	msgs := []Message{
		{Role: "user", Content: "What databases are currently active?"},
	}

	// TODO: In a real application, you would extract the user ID from the request context or session.
	actor := Actor{
		UserID:    uuid.Max,
		ProjectID: projectId,
	}

	result, err := h.client.Complete(r.Context(), msgs, actor)

	if err != nil {
		errormsg.WriteErrorResponse(w, r, err)
		return
	}
	render.JSON(w, r, result)

}
