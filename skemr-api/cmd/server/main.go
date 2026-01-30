package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lmittmann/tint"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/routers"
	"github.com/walmaa/skemr-api/internal/service"
	"github.com/walmaa/skemr-api/internal/tasks"
	"github.com/walmaa/skemr-api/internal/worker"
	"golang.org/x/net/context"
)

// runSchema drops the current schema, reads schema.sql file and executes it to set up the database schema.
func runSchema(conn *pgxpool.Pool) {
	schema, err := os.ReadFile("./db/schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Exec(context.Background(), "DROP SCHEMA IF EXISTS public CASCADE")
	_, err = conn.Exec(context.Background(), string(schema))
	if err != nil {
		log.Fatal(err)
	}
	// Seed data
	seed, err := os.ReadFile("./db/seed.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Exec(context.Background(), string(seed))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	ctx := context.Background()

	// Logger colors
	w := os.Stderr
	// Set global logger with custom options
	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.DateTime,
		}),
	))

	conn, err := pgxpool.New(context.Background(), "postgres://postgres:pass@localhost:5432/postgres")
	if err != nil {
		log.Fatal(err)
	}

	taskClient := tasks.StartTaskClient(ctx, "localhost:6379")
	queries := sqlc.New(conn)
	projectService := service.NewProjectService(queries)
	databaseService := service.NewDatabaseService(queries, taskClient)
	webhookService := service.NewWebhookService(queries)
	projectSecretsService := service.NewProjectSecretsService(queries)
	ruleService := service.NewRuleService(queries)
	databaseEntityService := service.NewDatabaseEntityService(queries)

	runSchema(conn)

	worker.StartTaskWorkers(queries)

	// Initialize services
	services := &routers.Services{
		ProjectService:        projectService,
		ProjectSecretsService: projectSecretsService,
		DatabaseService:       databaseService,
		WebhookService:        webhookService,
		RuleService:           ruleService,
		DatabaseEntityService: databaseEntityService,
	}

	// Initialize router
	router := routers.InitRouter(services)

	defer conn.Close()

	host := ":8080"
	srv := &http.Server{
		Addr:    host,
		Handler: router,
	}
	log.Printf("Listening and serving HTTP on %s", host)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
