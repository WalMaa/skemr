package ai

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/domain/databases"
	"github.com/walmaa/skemr-api/internal/domain/entities"
	"github.com/walmaa/skemr-common/models"
)

type ToolService struct {
	db sqlc.Querier
	connectorFactory func(database models.Database) DatabaseReader
}

func NewToolService(db sqlc.Querier) *ToolService {
	return &ToolService{db: db}
}

func (s *ToolService) GetDatabaseEntities(ctx context.Context, databaseId uuid.UUID, actor Actor) ([]models.DatabaseEntity, error) {
	slog.Info("ToolService: GetDatabaseEntities", "actor", actor, "databaseId", databaseId)

	databaseEntities, err := s.db.GetDatabaseEntitiesByDatabaseId(ctx, databaseId)

	if err != nil {
		slog.Error("Error getting database entities", "error", err)
		return nil, err
	}

	return entities.ToDomainDatabaseEntities(databaseEntities), nil

}

func (s *ToolService) GetDatabases(ctx context.Context, actor Actor) ([]models.Database, error) {
	slog.Info("ToolService: GetDatabases", "actor", actor)

	databaseRows, err := s.db.GetDatabasesByProjectId(ctx, actor.ProjectID)
	if err != nil {
		slog.Error("Error getting databases", "error", err)
		return nil, err
	}

	return databases.ToDomainDatabases(databaseRows), nil
}

func (s *ToolService) queryDatabase(ctx context.Context, databaseId uuid.UUID, query string, actor Actor) (string, error) {
	slog.Info("ToolService: queryDatabase", "actor", actor, "databaseId", databaseId, "query", query)



}