package models

import (
	"time"

	"github.com/google/uuid"
)

type DatabaseEntityType string

type DatabaseEntityStatus string

const (
	DatabaseEntityStatusActive  DatabaseEntityStatus = "active"
	DatabaseEntityStatusDeleted DatabaseEntityStatus = "deleted"
)

const (
	DatabaseEntityTypeDatabase DatabaseEntityType = "database"
	DatabaseEntityTypeSchema   DatabaseEntityType = "schema"
	DatabaseEntityTypeTable    DatabaseEntityType = "table"
	DatabaseEntityTypeColumn   DatabaseEntityType = "column"
)

type DatabaseEntity struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"` // Name of the entity "public", "users", "email", "my_view"
	Type        DatabaseEntityType     `json:"type"`
	ParentId    *uuid.UUID             `json:"parentId"` // in case of column, references table. table references schema etc.
	Status      DatabaseEntityStatus   `json:"status"`
	CreatedAt   time.Time              `json:"createdAt"`
	DeletedAt   *time.Time             `json:"deletedAt"`
	FirstSeenAt time.Time              `json:"firstSeenAt"` // when was this entity first seen in the database, useful for tracking new entities
	Attributes  map[string]interface{} `json:"attributes"`  // additional attributes like data type for column, nullability, etc.
}
