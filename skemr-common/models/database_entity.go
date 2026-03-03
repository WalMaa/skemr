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
	ID         uuid.UUID              `json:"id"`
	Name       string                 `json:"name"` // Name of the entity "public", "users", "email", "my_view"
	Type       DatabaseEntityType     `json:"type"`
	ParentId   *uuid.UUID             `json:"parentId"` // in case of column, references table. table references schema etc.
	ProjectId  uuid.UUID              `json:"projectId"`
	CreatedAt  time.Time              `json:"createdAt"`
	Attributes map[string]interface{} `json:"attributes"` // additional attributes like data type for column, nullability, etc.
}
