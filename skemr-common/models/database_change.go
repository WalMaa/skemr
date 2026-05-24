package models

import (
	"time"

	"github.com/google/uuid"
)

type DatabaseChange struct {
	Id         uuid.UUID                `json:"id"`
	DatabaseId uuid.UUID                `json:"databaseId"`
	EntityId   uuid.UUID                `json:"entityId"`
	Action     MigrationStatementAction `json:"action"`
	CreatedAt  time.Time                `json:"createdAt"`
}
