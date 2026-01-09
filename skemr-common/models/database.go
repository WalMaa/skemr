package models

import (
	"github.com/google/uuid"
)

type Database struct {
	ID           uuid.UUID    `json:"id"`
	DisplayName  string       `json:"displayName"`
	DbName       *string      `json:"dbName"`
	Username     *string      `json:"username"`
	Password     *string      `json:"password"`
	Host         *string      `json:"host"`
	Port         int32        `json:"port"`
	DatabaseType DatabaseType `json:"databaseType"`
	ProjectID    uuid.UUID    `json:"projectId"`
}

type DatabaseType string

const (
	Postgres DatabaseType = "postgres"
)
