package tasks

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

// A list of task types.
const (
	TypeDatabaseSync = "db:sync"
)

type DatabaseSyncPayload struct {
	DatabaseID uuid.UUID
}

func NewDatabaseSyncTask(databaseID uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(DatabaseSyncPayload{DatabaseID: databaseID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeDatabaseSync, payload, asynq.MaxRetry(3), asynq.Timeout(5*time.Minute)), nil
}
