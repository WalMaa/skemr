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
	r.Use(middleware.ErrorHandler())

	/*--------------------------- Public routes, auth, webhooks ---------------*/
	public := r.Group("/api/v1")

	auth := public.Group("/auth")
	{
		auth.POST("/login", nil)
		auth.POST("/register", nil)
	}

	// Webhook routes
	webhookController := controller.NewWebhookController()
	webhookController.RegisterRoutes(public)

	/*--------------------------- Protected routes ---------------------------*/
	protected := r.Group("/api/v1")

	// Auth
	protected.Use(gin.BasicAuth(gin.Accounts{
		"user": "pass",
	}))
	r.Use(middleware.AuthMiddleware())

	// Project routes, which will be prefixed with /api/v1/projects/:id
	// Each project will have its own set of routes under this group

	// Initialize controllers
	projectController := controller.NewProjectController(services.ProjectService)
	databaseController := controller.NewDatabaseController(services.DatabaseService)
	schemaController := controller.NewSchemaController(services.SchemaService)
	//ruleController := controller.NewRuleController(ruleService)

	// Register routes
	projectController.RegisterRoutes(protected)
	databaseController.RegisterRoutes(protected)
	schemaController.RegisterRoutes(protected)
	//ruleController.RegisterRoutes(projectRoutes)

	return r
}
