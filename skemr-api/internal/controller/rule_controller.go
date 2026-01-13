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
	group.GET("/:ruleId", h.GetRule)
	group.DELETE("/:ruleId", h.deleteRule)
}

func (h *RuleController) GetRule(c *gin.Context) {
	projectID := c.MustGet(middleware.CtxProjectID).(uuid.UUID)
	databaseId, ok := paramUUID(c, "databaseId")
	if !ok {
		return
	}
	ruleId, ok := paramUUID(c, "ruleId")
	if !ok {
		return
	}
	rule, err := h.Service.GetRule(c, projectID, databaseId, ruleId)

	if err != nil {
		c.Error(errors.New(err.Error()))
		return
	}

	c.JSON(http.StatusOK, rule)
}

func (h *RuleController) deleteRule(c *gin.Context) {
	projectID := c.MustGet(middleware.CtxProjectID).(uuid.UUID)
	databaseId, ok := paramUUID(c, "databaseId")
	if !ok {
		return
	}
	ruleId, ok := paramUUID(c, "ruleId")
	if !ok {
		return
	}
	err := h.Service.DeleteRule(c, projectID, databaseId, ruleId)

	if err != nil {
		c.Error(errors.New(err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
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
