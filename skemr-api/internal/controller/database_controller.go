package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/dto"
	"github.com/walmaa/skemr-api/internal/service"
)

type DatabaseController struct {
	Service *service.DatabaseService
}

func NewDatabaseController(s *service.DatabaseService) *DatabaseController {
	return &DatabaseController{Service: s}
}

func (h *DatabaseController) RegisterRoutes(g *gin.RouterGroup) {
	group := g.Group("/projects/:projectId/databases")
	{
		group.POST("/", h.createDatabase)
		group.GET("/:databaseId", h.getDatabase)
		group.DELETE("/:databaseId", h.deleteDatabase)
		group.GET("/", h.listDatabasesByProject)
		group.PATCH("/:databaseId", h.updateDatabase)
	}
}

func (h *DatabaseController) deleteDatabase(c *gin.Context) {
	databaseId, err := uuid.Parse(c.Param("databaseId"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	err = h.Service.DeleteDatabase(c, databaseId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(204) // No Content
}

func (h *DatabaseController) listDatabasesByProject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid project ID format"})
		return
	}

	databases, err := h.Service.ListDatabasesByProject(c, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, databases)
}

func (h *DatabaseController) createDatabase(c *gin.Context) {
	projectId, err := uuid.Parse(c.Param("projectId"))
	var body struct {
		Name string `json:"name" binding:"required,min=3,max=50"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	args := sqlc.CreateDatabaseParams{
		DisplayName: body.Name,
		ProjectID:   projectId,
	}

	database, err := h.Service.CreateDatabase(c, args)
	if err != nil {
		c.Error(errors.New(err.Error()))
		return
	}

	c.JSON(201, database)
}

func (h *DatabaseController) updateDatabase(c *gin.Context) {
	projectId := c.MustGet("projectID").(uuid.UUID)
	databaseId, ok := paramUUID(c, "databaseId")
	if !ok {
		return
	}
	var body dto.DatabaseUpdateDto
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	database, err := h.Service.UpdateDatabase(c, projectId, databaseId, body)
	if err != nil {
		c.Error(errors.New(err.Error()))
		return
	}

	c.JSON(http.StatusOK, database)
}

func (h *DatabaseController) getDatabase(c *gin.Context) {

	databaseId, err := uuid.Parse(c.Param("databaseId"))

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	database, err := h.Service.GetDatabase(c, databaseId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, database)
}
