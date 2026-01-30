package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/walmaa/skemr-api/internal/service"
)

type WebhookController struct {
	Service *service.WebhookService
}

func NewWebhookController(s *service.WebhookService) *WebhookController {
	return &WebhookController{Service: s}
}

func (h *WebhookController) RegisterRoutes(r chi.Router) {
	r.Route("/webhooks", func(r chi.Router) {
		r.Post("/gitlab", h.handleGitLabWebhook)
	})
}

func (h *WebhookController) handleGitLabWebhook(w http.ResponseWriter, r *http.Request) {
	slog.Info("Received GitLab webhook", "remote_addr", r.RemoteAddr)

	if err := h.Service.HandleGitLabWebhook(nil, r); err != nil { // TODO: Refactor service to not require gin.Context
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "GitLab webhook received"}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
