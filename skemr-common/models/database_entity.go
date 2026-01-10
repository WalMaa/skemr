package models

import (
	"time"

	"github.com/google/uuid"
)

type DatabaseEntityType string

const (
	DatabaseEntityTypeDatabase DatabaseEntityType = "database"
	DatabaseEntityTypeSchema   DatabaseEntityType = "schema"
	DatabaseEntityTypeTable    DatabaseEntityType = "table"
	DatabaseEntityTypeColumn   DatabaseEntityType = "column"
)

type DatabaseEntity struct {
	ID        uuid.UUID          `json:"id"`
	ProjectId uuid.UUID          `json:"projectId"`
	Type      DatabaseEntityType `json:"type"`
	ParentId  *uuid.UUID         `json:"parentId"` // in case of column, references table. table references schema etc.
	Name      string             `json:"name"`     // Name of the entity "public", "users", "email", "my_view"
	CreatedAt time.Time          `json:"createdAt"`
}
