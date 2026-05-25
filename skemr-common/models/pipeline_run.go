package models

import (
	"time"

	"github.com/google/uuid"
)

type PipelineRun struct {
	ID          uuid.UUID       `json:"id"`
	Status      MigrationStatus `json:"status"`
	Environment string          `json:"environment"`
	StartedAt   time.Time       `json:"startedAt"`
	CompletedAt time.Time       `json:"completedAt"`
	CreatedAt   time.Time       `json:"createdAt"`
}
