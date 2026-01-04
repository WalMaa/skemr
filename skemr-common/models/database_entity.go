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
	ID        uuid.UUID
	ProjectId uuid.UUID
	Type      DatabaseEntityType
	ParentId  *uuid.UUID // in case of column, references table. table references schema etc.
	Name      string     // Name of the entity "public", "users", "email", "my_view"
	CreatedAt time.Time
}
