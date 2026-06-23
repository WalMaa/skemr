package ai

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/mapper"
	"github.com/walmaa/skemr-common/models"
)

type ToolService struct {
	db sqlc.Querier
}

func NewToolService(db sqlc.Querier) *ToolService {
	return &ToolService{db: db}
}

func (s *ToolService) GetDatabaseEntities(ctx context.Context, databaseId uuid.UUID, actor Actor) ([]models.DatabaseEntity, error) {
	slog.Info("ToolService: GetDatabaseEntities", "actor", actor)

	entities, err := s.db.GetDatabaseEntities(ctx, sqlc.GetDatabaseEntitiesParams{
		DatabaseID: databaseId,
	})

	if err != nil {
		slog.Error("Error getting database entities", "error", err)
		return nil, err
	}

	return mapper.ToDomainDatabaseEntities(entities), nil

}
