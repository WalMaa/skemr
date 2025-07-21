package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/walmaa/skemr/db/sqlc"
	"github.com/walmaa/skemr/internal/service"
)

type DatabaseController struct {
	Service *service.DatabaseService
}

func NewDatabaseController(s *service.DatabaseService) *DatabaseController {
	return &DatabaseController{Service: s}
}

func (h *DatabaseController) RegisterRoutes(g *gin.RouterGroup) {
	group := g.Group("projects/:projectId/databases")
	{
		group.POST("/", h.createDatabase)
		group.GET("/:id", h.getDatabase)
		group.DELETE("/:id", h.deleteDatabase)
		group.GET("/", h.listDatabasesByProject)
	}
}

func (h *DatabaseController) deleteDatabase(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	err = h.Service.DeleteDatabase(c, id)
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
	id, err := uuid.Parse(c.Param("projectId"))
	var body struct {
		Name string `json:"name" binding:"required,min=3,max=50"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	args := sqlc.CreateDatabaseParams{
		Name:      body.Name,
		ProjectID: id,
	}

	database, err := h.Service.CreateDatabase(c, args)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, database)
}

func (h *DatabaseController) getDatabase(c *gin.Context) {

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	database, err := h.Service.GetDatabase(c, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, database)
}
