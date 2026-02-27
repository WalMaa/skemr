package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/walmaa/skemr-api/internal/controller"
	"github.com/walmaa/skemr-api/internal/middleware"
	"github.com/walmaa/skemr-api/internal/service"
)

type Services struct {
	ProjectService        *service.ProjectService
	DatabaseService       *service.DatabaseService
	RuleService           *service.RuleService
	WebhookService        *service.WebhookService
	ProjectSecretsService *service.ProjectSecretsService
	DatabaseEntityService *service.DatabaseEntityService
}

func InitRouter(services *Services) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.StripSlashes)

	// Public routes
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/login", nil)
		r.Post("/register", nil)
	})

	webhookController := controller.NewWebhookController(services.WebhookService)
	webhookController.RegisterRoutes(r)

	// Protected routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		projectController := controller.NewProjectController(services.ProjectService)

		r.Get("/projects", projectController.GetProjects)
		r.Post("/projects", projectController.CreateProject)

		// Project level routes
		r.Route("/projects/{projectId}", func(r chi.Router) {
			r.Use(middleware.ProjectIDMiddleware)
			databaseController := controller.NewDatabaseController(services.DatabaseService)
			projectSecretsController := controller.NewProjectSecretsController(services.ProjectSecretsService)
			ruleController := controller.NewRuleController(services.RuleService)
			databaseEntityController := controller.NewDatabaseEntityController(services.DatabaseEntityService)
			databaseController.RegisterRoutes(r)
			projectSecretsController.RegisterRoutes(r)
			ruleController.RegisterRoutes(r)
			databaseEntityController.RegisterRoutes(r)

			r.Get("/", projectController.GetProject)
			r.Delete("/", projectController.DeleteProject)
		})
	})

	return r
}
