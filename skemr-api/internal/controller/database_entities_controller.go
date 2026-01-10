package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/internal/middleware"
	"github.com/walmaa/skemr-api/internal/service"
	"github.com/walmaa/skemr-common/models"
)

type DatabaseEntityController struct {
	Service *service.DatabaseEntityService
}

func NewDatabaseEntityController(s *service.DatabaseEntityService) *DatabaseEntityController {
	return &DatabaseEntityController{Service: s}
}

func (h *DatabaseEntityController) RegisterRoutes(g *gin.RouterGroup) {
	group := g.Group("/projects/:projectId/databases/:databaseId/entities")
	{
		group.GET("", h.GetDatabaseEntities)
		group.GET("/:entityId", h.GetDatabaseEntity)

	}
}

type Query struct {
	EntityType models.DatabaseEntityType
}

func (h *DatabaseEntityController) GetDatabaseEntity(c *gin.Context) {
	projectId := c.MustGet(middleware.CtxProjectID).(uuid.UUID)

	databaseId, err := uuid.Parse(c.Param("databaseId"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid database ID format"})
		return
	}
	entityId, err := uuid.Parse(c.Param("entityId"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid entity ID format"})
		return
	}
	entity, err := h.Service.GetDatabaseEntityByID(c, projectId, databaseId, entityId)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entity)
}

func (h *DatabaseEntityController) GetDatabaseEntities(c *gin.Context) {
	projectId := c.MustGet(middleware.CtxProjectID).(uuid.UUID)

	databaseId, err := uuid.Parse(c.Param("databaseId"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid database ID format"})
		return
	}
	entityTypeQuery := c.Query("type")
	var entityType *models.DatabaseEntityType
	if entityTypeQuery != "" {
		et := models.DatabaseEntityType(entityTypeQuery)
		entityType = &et
	}

	parentIdQuery := c.Query("parentId")
	var parentId *uuid.UUID
	if parentIdQuery != "" {
		pId, err := uuid.Parse(parentIdQuery)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid parent ID format"})
			return
		}
		parentId = &pId
	}
	entities, err := h.Service.ListDatabaseEntitiesByDatabase(c, projectId, databaseId, entityType, parentId)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entities)

}
