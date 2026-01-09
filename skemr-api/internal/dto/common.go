package dto

import "github.com/google/uuid"

type DatabaseCreationDto struct {
	DisplayName string `json:"display_name"`
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
