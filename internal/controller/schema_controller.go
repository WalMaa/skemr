package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/walmaa/skemr/db/sqlc"
	"github.com/walmaa/skemr/internal/service"
)

type SchemaController struct {
	Service *service.SchemaService
}

func NewSchemaController(s *service.SchemaService) *SchemaController {
	return &SchemaController{Service: s}
}

func (h *SchemaController) RegisterRoutes(g *gin.RouterGroup) {
	group := g.Group("projects/:projectId/schemas")
	group.POST("/", h.createSchema)
	group.GET("/:schemaId", h.getSchema)
	//group.PUT("/:schemaId", h.updateSchema)
	//group.DELETE("/:schemaId", h.deleteSchema)
}

func (h *SchemaController) getSchema(c *gin.Context) {

	id, err := uuid.Parse(c.Param("schemaId"))

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid schema ID format"})
		return
	}

	projectId, err := uuid.Parse(c.Param("projectId"))

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid project ID format"})
		return
	}

	schema, err := h.Service.GetSchema(c, projectId, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, schema)
}

func (h *SchemaController) createSchema(c *gin.Context) {
	var body struct {
		Name       string    `json:"name" binding:"required,min=3,max=50"`
		DatabaseID uuid.UUID `json:"database_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	projectId, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid project ID format"})
		return
	}

	schema, err := h.Service.CreateSchema(c, projectId, sqlc.CreateSchemaParams{
		Name:       body.Name,
		DatabaseID: body.DatabaseID,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, schema)
}
