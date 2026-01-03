package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/internal/service"
	"github.com/walmaa/skemr-common/models"
)

type RuleController struct {
	Service *service.RuleService
}

func NewRuleController(s *service.RuleService) *RuleController {
	return &RuleController{Service: s}
}

func (h *RuleController) RegisterRoutes(g *gin.RouterGroup) {
	group := g.Group("projects/:projectId/rules")
	group.POST("/", h.createRule)
}

func (h *RuleController) createRule(c *gin.Context) {
	projectId, err := uuid.Parse(c.Param("projectId"))

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var body models.RuleCreationDto
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	rule, err := h.Service.CreateRule(c, projectId, body)

	if err != nil {
		c.Error(errors.New(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, rule)
}
