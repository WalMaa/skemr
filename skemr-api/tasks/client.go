package tasks

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/hibiken/asynq"
)

func StartTaskClient(ctx context.Context, addr string) *asynq.Client {
	slog.Info("Starting Asynq client")
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: addr})

	go func() {
		<-ctx.Done()
		// optional tiny delay; usually unnecessary
		time.Sleep(10 * time.Millisecond)

		_ = client.Close()
		log.Println("asynq client stopped")
	}()

	return client
}
