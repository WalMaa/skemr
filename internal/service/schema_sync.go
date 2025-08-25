package service

import (
	"github.com/walmaa/skemr/db/sqlc"
	"github.com/walmaa/skemr/dbreflect"
)

type SchemaSyncService struct {
	connector *dbreflect.PostgresConnector
	db        sqlc.Querier
}

func NewSchemaSyncService(connector *dbreflect.PostgresConnector, db sqlc.Querier) *SchemaSyncService {
	return &SchemaSyncService{connector: connector, db: db}
}

func (s *SchemaSyncService) SyncSchema(databaseID, schemaID string) error {
	return nil
}
