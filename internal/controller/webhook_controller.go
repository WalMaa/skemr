package controller

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type WebhookController struct {
}

func NewWebhookController() *WebhookController {
	return &WebhookController{}
}

func (h *WebhookController) RegisterRoutes(g *gin.RouterGroup) {
	group := g.Group("webhooks")
	{
		group.POST("/gitlab", h.handleGitLabWebhook)

	}
}

func (h *WebhookController) handleGitLabWebhook(c *gin.Context) {
	// Handle GitLab webhook payload
	slog.Info("Received GitLab webhook", "body", c.Request)
	c.JSON(200, gin.H{"status": "GitLab webhook received"})
}
