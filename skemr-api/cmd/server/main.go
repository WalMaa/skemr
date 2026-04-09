package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/lmittmann/tint"
	"github.com/pressly/goose/v3"
	"github.com/walmaa/skemr-api/config"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/routers"
	"github.com/walmaa/skemr-api/internal/service"
	"github.com/walmaa/skemr-api/internal/tasks"
	"github.com/walmaa/skemr-api/internal/validation"
	"github.com/walmaa/skemr-api/internal/worker"
)

// runSchema sets up the database schema.
func runSchema(conn *pgxpool.Pool) {

	_, err := conn.Exec(context.Background(), `
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;
	`)

	if err != nil {
		slog.Error("Error refreshing schema", "error", err)
	}

	db := stdlib.OpenDBFromPool(conn)

	goose.SetDialect("postgres")
	if err := goose.Up(db, "./db/migrations"); err != nil {
		panic(err)
	}

}

func seedTestData(conn *pgxpool.Pool) {
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

	validation.Init()

	// Logger colors
	w := os.Stderr
	// Set global logger with custom options
	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.DateTime,
		}),
	))

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Error loading config", "error", err)
		os.Exit(1)
	}

	slog.Info("Starting Skemr API server", "environment", cfg.App.Env)

	slog.Info("Connecting to database", "host", cfg.Database.Host, "port", cfg.Database.Port, "name", cfg.Database.Name, "sslmode", cfg.Database.SSLMode)
	conn, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name))
	if err != nil {
		slog.Error("Error connecting to DB", "error", err)
		os.Exit(1)
	}

	taskClient := tasks.StartTaskClient(ctx, asynq.RedisClientOpt{
		Addr:      fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password:  cfg.Redis.Password,
		DB:        cfg.Redis.DB,
		TLSConfig: nil,
	})

	queries := sqlc.New(conn)
	projectService := service.NewProjectService(queries)
	databaseService := service.NewDatabaseService(queries, taskClient)
	webhookService := service.NewWebhookService(queries)
	projectSecretsService := service.NewAccessTokenService(queries)
	ruleService := service.NewRuleService(queries)
	databaseEntityService := service.NewDatabaseEntityService(queries)
	integrationService := service.NewIntegrationService(ruleService)

	if cfg.App.Env == "dev" {
		runSchema(conn)
		seedTestData(conn)
	}

	worker.StartTaskWorkers(queries, cfg)

	// Initialize services
	services := &routers.Services{
		ProjectService:        projectService,
		AccessTokenService:    projectSecretsService,
		DatabaseService:       databaseService,
		WebhookService:        webhookService,
		RuleService:           ruleService,
		DatabaseEntityService: databaseEntityService,
		IntegrationService:    integrationService,
	}

	// Initialize router
	router := routers.InitRouter(services)

	defer conn.Close()

	host := fmt.Sprintf(":%d", cfg.App.Port)
	srv := &http.Server{
		Addr:    host,
		Handler: router,
	}
	slog.Info("Listening and serving HTTP", "host", host)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("HTTP server ListenAndServe error", "error", err)
		os.Exit(1)
	}
}
