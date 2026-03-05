package worker

import (
	"log"
	"log/slog"

	"github.com/hibiken/asynq"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/dbreflect"
	"github.com/walmaa/skemr-api/internal/tasks"
)

func StartTaskWorkers(db sqlc.Querier) {

	srv := asynq.NewServer(asynq.RedisClientOpt{Addr: "localhost:6379"}, asynq.Config{
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
