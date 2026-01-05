package models

import (
	"github.com/google/uuid"
)

type MigrationStatementAction string

const (
	MigrationStatementActionCreate MigrationStatementAction = "create"
	MigrationStatementActionAlter  MigrationStatementAction = "alter"
	MigrationStatementActionDrop   MigrationStatementAction = "drop"
	MigrationStatementActionInsert MigrationStatementAction = "insert"
	MigrationStatementActionUpdate MigrationStatementAction = "update"
	MigrationStatementActionDelete MigrationStatementAction = "delete"
)

type MigrationStatus string

const (
	MigrationStatusPending    MigrationStatus = "pending"
	MigrationStatusInProgress MigrationStatus = "in_progress"
	MigrationStatusCompleted  MigrationStatus = "completed"
	MigrationStatusFailed     MigrationStatus = "failed"
)

type MigrationStatement struct {
	ID           uuid.UUID                `json:"id"`
	SchemaID     uuid.UUID                `json:"schema_id"`
	RawStatement string                   `json:"raw_statement"`
	Action       MigrationStatementAction `json:"action"`
	Status       MigrationStatus          `json:"status"`
	Target       *string                  `json:"target"`
	RelationName *string                  `json:"relation_name"`
}
