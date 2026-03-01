package dto

import (
	"github.com/google/uuid"
)

type ProjectCreationDto struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}

type DatabaseCreationDto struct {
	DisplayName  string       `json:"displayName" validate:"required"`
	DbName       *string      `json:"dbName"`
	Username     *string      `json:"username"`
	Password     *string      `json:"password"`
	Host         *string      `json:"host"`
	Port         int32        `json:"port"`
	DatabaseType DatabaseType `json:"databaseType" validate:"required,oneof=postgres"`
}

type DatabaseUpdateDto struct {
	DisplayName  *string       `json:"displayName"`
	DbName       *string       `json:"dbName"`
	Username     *string       `json:"username"`
	Password     *string       `json:"password"`
	Host         *string       `json:"host"`
	Port         *int32        `json:"port"`
	DatabaseType *DatabaseType `json:"databaseType"`
}

type DatabaseType string

const (
	Postgres DatabaseType = "postgres"
)

type RuleCreationDto struct {
	Name             string
	RuleType         RuleType
	DataBaseEntityId uuid.UUID
}

type RuleType string

const (
	RuleTypeLocked     RuleType = "locked"
	RuleTypeWarn       RuleType = "warn"
	RuleTypeAdvisory   RuleType = "advisory"
	RuleTypeDeprecated RuleType = "deprecated"
)

type SecretCreationDto struct {
	Name      string `json:"name" validate:"required"`
	ExpiresAt string `json:"expiresAt" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}
