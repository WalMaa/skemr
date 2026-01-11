package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/internal/dto"
	"github.com/walmaa/skemr-api/internal/middleware"
	"github.com/walmaa/skemr-api/internal/service"
)

type RuleController struct {
	Service *service.RuleService
}

func NewRuleController(s *service.RuleService) *RuleController {
	return &RuleController{Service: s}
}

func (h *RuleController) RegisterRoutes(g *gin.RouterGroup) {
	group := g.Group("projects/:projectId/databases/:databaseId/rules")
	group.POST("", h.createRule)
	group.GET("", h.ListRules)
}

func (h *RuleController) createRule(c *gin.Context) {
	databaseId, ok := paramUUID(c, "databaseId")
	if !ok {
		return
	}
	projectID := c.MustGet(middleware.CtxProjectID).(uuid.UUID)

	var body dto.RuleCreationDto
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	rule, err := h.Service.CreateRule(c, projectID, databaseId, body)

	if err != nil {
		c.Error(errors.New(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, rule)
}

func (h *RuleController) ListRules(c *gin.Context) {
	databaseId, ok := paramUUID(c, "databaseId")
	if !ok {
		return
	}
	projectID := c.MustGet("projectID").(uuid.UUID)

	rules, err := h.Service.ListRulesByDatabase(c, projectID, databaseId)

	if err != nil {
		c.Error(errors.New(err.Error()))
		return
	}

	c.JSON(http.StatusOK, rules)

}
