package controller

import (
	"io"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/walmaa/skemr/internal/service"
)

type WebhookController struct {
	Service *service.WebhookService
}

func NewWebhookController(s *service.WebhookService) *WebhookController {
	return &WebhookController{Service: s}
}

func (h *WebhookController) RegisterRoutes(g *gin.RouterGroup) {
	group := g.Group("webhooks")
	{
		group.POST("/gitlab", h.handleGitLabWebhook)

	}
}

func (h *WebhookController) handleGitLabWebhook(c *gin.Context) {
	// Handle GitLab webhook payload
	slog.Info("Received GitLab webhook", "body", c.Request.RemoteAddr)
	webhookSecret := c.GetHeader("X-Gitlab-Token")
	if webhookSecret == "" {
		slog.Warn("Missing X-Gitlab-Token header")
		c.JSON(400, gin.H{"error": "Missing X-Gitlab-Token header"})
		return
	}
	jsonData, err := io.ReadAll(c.Request.Body)

	if err != nil {
		slog.Error("Failed to read request body", "error", err)
		c.JSON(400, gin.H{"error": "Failed to read request body"})
		return
	}
	h.Service.HandleGitLabWebhook(c, jsonData, webhookSecret)

	c.JSON(200, gin.H{"status": "GitLab webhook received"})

}
