package controller

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/walmaa/skemr-api/internal/service"
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

	r := c.Request

	h.Service.HandleGitLabWebhook(c, r)

	c.JSON(200, gin.H{"status": "GitLab webhook received"})

}
