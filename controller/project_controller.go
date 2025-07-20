package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"skemr/service"
)

type ProjectController struct {
	Service *service.ProjectService
}

func NewProjectController(s *service.ProjectService) *ProjectController {
	return &ProjectController{Service: s}
}

func (h *ProjectController) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/projects")
	group.POST("/", h.createProject)
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
