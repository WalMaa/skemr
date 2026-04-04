package worker

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/hibiken/asynq"
	"github.com/walmaa/skemr-api/config"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/dbreflect"
	"github.com/walmaa/skemr-api/internal/tasks"
)

func StartTaskWorkers(db sqlc.Querier, cfg *config.Config) {

	srv := asynq.NewServer(asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)}, asynq.Config{
		Concurrency: 10,
	})
	syncService := dbreflect.NewSchemaSyncService(db, dbreflect.NewPostgresConnector)

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeDatabaseSync, syncService.ProcessSyncTask)

	go func() {
		slog.Info("asynq worker starting...")
		if err := srv.Run(mux); err != nil {
			log.Fatalf("asynq worker stopped: %v", err)
		}
	}()
}
