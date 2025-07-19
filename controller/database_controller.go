package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	skemr "skemr/db/sqlc"
	"skemr/service"
)

type DatabaseController struct {
	Service *service.DatabaseService
}

func NewDatabaseController(s *service.DatabaseService) *DatabaseController {
	return &DatabaseController{Service: s}
}

func (h *DatabaseController) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/databases")
	{
		group.POST("/", h.createDatabase)
		group.GET("/:id", h.getDatabase)
		group.DELETE("/:id", h.deleteDatabase)
		group.GET("/project/:id", h.listDatabasesByProject)
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
	id, err := uuid.Parse(c.Param("id"))
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
	var body struct {
		Name     string
		Username string
		Password string
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	args := skemr.CreateDatabaseParams{
		Name:     body.Name,
		Username: &body.Username,
		Password: &body.Password,
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
