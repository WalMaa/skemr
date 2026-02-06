package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/internal/dto"
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
		r.Post("/", h.createToken)
		r.Get("/", h.getSecrets)
		r.Get("/{secretId}", h.getSecret)
		r.Put("/{secretId}", h.updateSecret)
		r.Delete("/{secretId}", h.deleteSecret)
	})
}

func (h *ProjectSecretsController) createToken(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	projectId := c.Value("projectId").(uuid.UUID)
	var req dto.SecretCreationDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	token, err := h.Service.CreateToken(c, projectId, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type TokenResponse struct {
		Token string `json:"token"`
	}
	response := TokenResponse{Token: token}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, response)

}

func (h *ProjectSecretsController) getSecret(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement retrieval logic
}

func (h *ProjectSecretsController) getSecrets(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	projectId := c.Value("projectId").(uuid.UUID)

	tokens, err := h.Service.GetTokens(c, projectId)
	if err != nil {
		slog.Error("Error getting tokens", err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, tokens)
}

func (h *ProjectSecretsController) updateSecret(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement update logic
}

func (h *ProjectSecretsController) deleteSecret(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	projectId := c.Value("projectId").(uuid.UUID)

	secretId, err := uuid.Parse(chi.URLParam(r, "secretId"))

	if err != nil {
		http.Error(w, "Invalid Secret ID", http.StatusBadRequest)
	}

	h.Service.DeleteToken(c, projectId, secretId)

	render.Status(r, http.StatusNoContent)
}
