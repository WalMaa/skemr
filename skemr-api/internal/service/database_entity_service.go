package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/mapper"
	"github.com/walmaa/skemr-common/models"
)

type DatabaseEntityService struct {
	db sqlc.Querier
}

func NewDatabaseEntityService(q sqlc.Querier) *DatabaseEntityService {
	{
		return &DatabaseEntityService{db: q}
	}
}

func (s *DatabaseEntityService) GetDatabaseEntityByID(c context.Context, projectId uuid.UUID, databaseId uuid.UUID, entityId uuid.UUID) (models.DatabaseEntity, error) {
	slog.Info("Getting database entity", "projectId", projectId, "database", databaseId, "entityId", entityId)
	project, err := CheckProjectExists(c, s.db, projectId)

	if err != nil {
		return models.DatabaseEntity{}, err
	}

	database, err := CheckDatabaseExists(c, s.db, project.ID, databaseId)
	if err != nil {
		return models.DatabaseEntity{}, err
	}

	entity, err := s.db.GetDatabaseEntity(c, database.ID)

	return mapper.ToDomainDatabaseEntity(entity), err
}

func (s *DatabaseEntityService) ListDatabaseEntitiesByDatabase(c context.Context, projectId uuid.UUID, databaseId uuid.UUID, entityType *models.DatabaseEntityType, parentId *uuid.UUID) ([]models.DatabaseEntity, error) {
	slog.Info("Listing database entities", "projectId", projectId, "database", databaseId)
	project, err := CheckProjectExists(c, s.db, projectId)

	if err != nil {
		return []models.DatabaseEntity{}, err
	}

	database, err := CheckDatabaseExists(c, s.db, project.ID, databaseId)
	if err != nil {
		return []models.DatabaseEntity{}, err
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

	return mapper.ToDomainDatabaseEntities(entities), nil
}
