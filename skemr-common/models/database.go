package models

import (
	"github.com/google/uuid"
)

type Database struct {
	ID          uuid.UUID    `json:"id"`
	DisplayName string       `json:"display_name"`
	DbName      string       `json:"db_name"`
	Username    *string      `json:"username"`
	Password    *string      `json:"password"`
	Host        *string      `json:"host"`
	Port        int32        `json:"port"`
	Type        DatabaseType `json:"type"`
	ProjectID   uuid.UUID    `json:"project_id"`
}

type DatabaseType string

const (
	Postgres DatabaseType = "postgres"
)

type DatabaseCreationDto struct {
	DisplayName string `json:"display_name"`
}

type DatabaseReturnDto struct {
	ID          uuid.UUID    `json:"id"`
	DisplayName string       `json:"display_name"`
	DbName      string       `json:"db_name"`
	Username    *string      `json:"username"`
	Host        *string      `json:"host"`
	Port        int32        `json:"port"`
	Type        DatabaseType `json:"type"`
	ProjectID   uuid.UUID    `json:"project_id"`
}
