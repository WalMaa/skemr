package entities

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/service"
	"github.com/walmaa/skemr-common/models"
)

type DatabaseEntityService struct {
	db            sqlc.Querier
	scopeResolver service.ScopeResolver
}

func NewDatabaseEntityService(q sqlc.Querier, scopeResolver service.ScopeResolver) *DatabaseEntityService {
	{
		return &DatabaseEntityService{db: q, scopeResolver: scopeResolver}
	}
}

func (s *DatabaseEntityService) GetDatabaseEntityByID(c context.Context, projectId uuid.UUID, databaseId uuid.UUID, entityId uuid.UUID) (models.DatabaseEntity, error) {
	slog.Info("Getting database entity", "projectId", projectId, "database", databaseId, "entityId", entityId)

	_, err := s.scopeResolver.RequireDatabase(c, projectId, databaseId)
	if err != nil {
		slog.Error("Error getting database", "err", err)
		return models.DatabaseEntity{}, err
	}

	entity, err := s.db.GetDatabaseEntity(c, entityId)

	return ToDomainDatabaseEntity(entity), err
}

func (s *DatabaseEntityService) ListDatabaseEntitiesByDatabase(c context.Context, projectId uuid.UUID, databaseId uuid.UUID, entityType *models.DatabaseEntityType, parentId *uuid.UUID) ([]models.DatabaseEntity, error) {
	slog.Info("Listing database entities", "projectId", projectId, "database", databaseId)
	database, err := s.scopeResolver.RequireDatabase(c, projectId, databaseId)
	if err != nil {
		slog.Error("Error getting database", "err", err)
		return nil, err
	}

	var et sqlc.NullDatabaseEntityType
	if entityType != nil {
		et = sqlc.NullDatabaseEntityType{
			DatabaseEntityType: sqlc.DatabaseEntityType(*entityType),
			Valid:              true,
		}
	}

	entities, err := s.db.GetDatabaseEntities(c, sqlc.GetDatabaseEntitiesParams{
		DatabaseID: database.ID,
		EntityType: et,
		ParentID:   parentId,
	})
	if err != nil {
		return []models.DatabaseEntity{}, err
	}

	return ToDomainDatabaseEntities(entities), nil
}
