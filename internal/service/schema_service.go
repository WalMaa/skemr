package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/walmaa/skemr/db/sqlc"
	"log/slog"
)

type SchemaService struct {
	db sqlc.Querier
}

func NewSchemaService(q sqlc.Querier) *SchemaService {
	return &SchemaService{db: q}
}

func (s *SchemaService) CreateSchema(c context.Context, projectID uuid.UUID, args sqlc.CreateSchemaParams) (sqlc.Schema, error) {
	slog.Info("Creating schema", "name", args)

	// Check if the project exists
	database, err := CheckDatabaseExists(c, s.db, projectID, args.DatabaseID)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("Error checking if project exists", "project_id", projectID, "err", err)
		return sqlc.Schema{}, err
	}

	// Check if a schema with the given name already exists
	exists, err := s.db.GetSchemaByNameAndDatabase(c, sqlc.GetSchemaByNameAndDatabaseParams{
		Name:       args.Name,
		DatabaseID: database.ID,
	},
	)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("Error checking for existing schema", "name", args.Name, "database_id", database.ID, "err", err)
		return sqlc.Schema{}, err
	}

	// If exists is not empty, it means a schema with the same name already exists.
	if exists != (sqlc.Schema{}) {
		slog.Warn("Schema already exists", "name", args.Name, "database_id", database.ID)
		return sqlc.Schema{}, errors.New("schema already exists")
	}

	// Create the schema
	return s.db.CreateSchema(c, args)
}

func (s *SchemaService) GetSchema(c context.Context, projectID uuid.UUID, schemaID uuid.UUID) (sqlc.Schema, error) {
	slog.Info("Getting schema", "schema_id", schemaID)

	// Check if the project exists
	project, err := CheckProjectExists(c, s.db, projectID)
	if err != nil {
		slog.Error("Error checking if project exists", "project_id", projectID, "err", err)
		return sqlc.Schema{}, err
	}

	// Check if the schema exists
	schema, err := s.db.GetSchemaByIdAndProject(c, sqlc.GetSchemaByIdAndProjectParams{
		ID:        schemaID,
		ProjectID: project.ID,
	})
	if err != nil {
		slog.Error("Error getting schema", "schema_id", schemaID, "err", err)
		return sqlc.Schema{}, err
	}

	return schema, nil
}
