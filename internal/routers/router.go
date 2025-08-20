package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/walmaa/skemr/internal/controller"
	"github.com/walmaa/skemr/internal/middleware"
	"github.com/walmaa/skemr/internal/service"
)

type Services struct {
	ProjectService  *service.ProjectService
	DatabaseService *service.DatabaseService
	SchemaService   *service.SchemaService
}

func InitRouter(services *Services) *gin.Engine {
	r := gin.Default()

	// Auth
	r.Use(gin.BasicAuth(gin.Accounts{
		"user": "pass",
	}))
	r.Use(middleware.AuthMiddleware())

	r.Use(middleware.ErrorHandler())

	auth := r.Group("/auth")
	{
		auth.POST("/login", nil)
		auth.POST("/register", nil)
	}

	// API routes
	api := r.Group("/api")
	// Project routes, which will be prefixed with /api/projects/:id
	// Each project will have its own set of routes under this group

	// Initialize controllers
	projectController := controller.NewProjectController(services.ProjectService)
	databaseController := controller.NewDatabaseController(services.DatabaseService)
	schemaController := controller.NewSchemaController(services.SchemaService)
	//ruleController := controller.NewRuleController(ruleService)

	// Register routes
	projectController.RegisterRoutes(api)
	databaseController.RegisterRoutes(api)
	schemaController.RegisterRoutes(api)
	//ruleController.RegisterRoutes(projectRoutes)

	return r
}
