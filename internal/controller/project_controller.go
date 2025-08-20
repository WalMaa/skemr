package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	errormsg "github.com/walmaa/skemr/errormsg"
	"github.com/walmaa/skemr/internal/service"
)

type ProjectController struct {
	Service *service.ProjectService
}

func NewProjectController(s *service.ProjectService) *ProjectController {
	return &ProjectController{Service: s}
}

func (h *ProjectController) RegisterRoutes(g *gin.RouterGroup) {
	group := g.Group("/projects/")
	group.POST("/", h.createProject)
	group.GET("/:projectId", h.getProject)
	group.GET("/", h.getProjects)
}

func (h *ProjectController) getProjects(c *gin.Context) {
	projects, err := h.Service.GetProjects(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}

func (h *ProjectController) createProject(c *gin.Context) {
	var body struct {
		Name string `json:"name" binding:"required,min=3,max=50"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.CreateProject(c, body.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *ProjectController) getProject(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		c.Error(errormsg.ErrInvalidIdFormat)
		return
	}

	project, err := h.Service.GetProject(c, projectID)
	if err != nil {
		c.Error(errors.New(err.Error()))
		return
	}

	c.JSON(http.StatusOK, project)
}
