package service

import "github.com/walmaa/skemr/db/sqlc"

type SchemaService struct {
	db sqlc.Querier
}

func NewSchemaService(q sqlc.Querier) *SchemaService {
	return &SchemaService{db: q}
}

func (s *SchemaService) CreateSchema(args sqlc.CreateSchemaParams) (sqlc.Schema, error) {

}
