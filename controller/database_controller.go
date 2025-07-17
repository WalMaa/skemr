package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	skemr "skemr/db"
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
		Username: pgtype.Text.Scan(body.Username),
		Password: body.Password,
	}

	database, err := h.Service.CreateDatabase(c, args)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, database)
}

func (h *DatabaseController) getDatabase(c *gin.Context) {
	id, err := pgtype.UUID{}.Scan(c.Param("id"))
	pgtype.UUID{}
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	database, err := h.Service.GetDatabase(c, id.(pgtype.UUID))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, database)
}
