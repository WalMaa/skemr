package routers

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"skemr/controller"
	"skemr/docs"
	"skemr/service"
)

type Services struct {
	ProjectService  *service.ProjectService
	DatabaseService *service.DatabaseService
}

func InitRouter(services *Services) *gin.Engine {
	r := gin.Default()

	// Auth
	r.Use(gin.BasicAuth(gin.Accounts{
		"user": "pass",
	}))

	auth := r.Group("/auth")
	{
		auth.POST("/login", nil)
		auth.POST("/register", nil)
	}

	// Swagger documentation
	docs.SwaggerInfo.BasePath = "/api"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// API routes
	api := r.Group("/api")
	// Project routes, which will be prefixed with /api/projects/:id
	// Each project will have its own set of routes under this group

	// Initialize controllers
	projectController := controller.NewProjectController(services.ProjectService)
	databaseController := controller.NewDatabaseController(services.DatabaseService)
	//ruleController := controller.NewRuleController(ruleService)

	// Register routes
	projectController.RegisterRoutes(api)
	databaseController.RegisterRoutes(api)
	//ruleController.RegisterRoutes(projectRoutes)

	return r
}
