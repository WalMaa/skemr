package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/walmaa/skemr-api/internal/service"
)

type ProjectSecretsController struct {
	Service *service.ProjectSecretsService
}

func NewProjectSecretsController(s *service.ProjectSecretsService) *ProjectSecretsController {
	return &ProjectSecretsController{Service: s}
}

func (h *ProjectSecretsController) RegisterRoutes(r chi.Router) {
	r.Route("/secrets", func(r chi.Router) {
		r.Post("/", h.createSecret)
		r.Get("/", h.getSecrets)
		r.Get("/{secretId}", h.getSecret)
		r.Put("/{secretId}", h.updateSecret)
		r.Delete("/{secretId}", h.deleteSecret)
	})
}

func (h *ProjectSecretsController) createSecret(w http.ResponseWriter, r *http.Request) {
	type createSecretRequest struct {
		Name      string  `json:"name"`
		ExpiresAt *string `json:"expiresAt"`
	}
	var req createSecretRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	// TODO: Implement creation logic
}

func (h *ProjectSecretsController) getSecret(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement retrieval logic
}

func (h *ProjectSecretsController) getSecrets(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement listing logic
}

func (h *ProjectSecretsController) updateSecret(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement update logic
}

func (h *ProjectSecretsController) deleteSecret(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement delete logic
}
