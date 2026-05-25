package models

import (
	"time"

	"github.com/google/uuid"
)

type DatabaseChange struct {
	Id        uuid.UUID                `json:"id"`
	EntityId  uuid.UUID                `json:"entityId"`
	Action    MigrationStatementAction `json:"action"`
	CreatedAt time.Time                `json:"createdAt"`
}
