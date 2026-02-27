package tasks

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/hibiken/asynq"
)

func StartTaskClient(ctx context.Context, clientOpt asynq.RedisClientOpt) *asynq.Client {
	slog.Info("Starting Asynq client", "redis_addr", clientOpt.Addr, "redis_db", clientOpt.DB)
	client := asynq.NewClient(clientOpt)

	go func() {
		<-ctx.Done()
		// optional tiny delay; usually unnecessary
		time.Sleep(10 * time.Millisecond)

		_ = client.Close()
		log.Println("asynq client stopped")
	}()

	return client
}
