package service

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/errormsg"
	"github.com/walmaa/skemr-api/internal/mapper"
	"github.com/walmaa/skemr-common/models"
)

type ScopeResolver interface {
	RequireDatabase(c context.Context, projectId uuid.UUID, databaseId uuid.UUID) (models.Database, error)
	RequireDatabaseEntity(c context.Context, projectId uuid.UUID, databaseId uuid.UUID, entityId uuid.UUID) (models.DatabaseEntity, error)
}
type SqlcScopeResolver struct {
	db sqlc.Querier
}

func NewScopeResolver(q sqlc.Querier) *SqlcScopeResolver {
	return &SqlcScopeResolver{db: q}
}

func (s *SqlcScopeResolver) RequireDatabase(c context.Context, projectId uuid.UUID, databaseId uuid.UUID) (models.Database, error) {
	slog.Info("Getting database", "databaseId", databaseId, "projectId", projectId)

	database, err := s.db.GetDatabaseByIDAndProjectID(c, sqlc.GetDatabaseByIDAndProjectIDParams{
		ID:        databaseId,
		ProjectID: projectId,
	})

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Info("Database not found", "databaseId", databaseId, "projectId", projectId)
		return models.Database{}, &models.ErrorResponse{
			Message: errormsg.ErrDatabaseNotFound,
			Errors:  nil,
			Status:  http.StatusBadRequest,
		}
	}

	if err != nil {
		slog.Error("Unable to get database", "databaseId", databaseId, "projectId", projectId, "err", err)
		return models.Database{}, err
	}

	return mapper.ToDomainDatabase(database), nil
}

func (s *SqlcScopeResolver) RequireDatabaseEntity(c context.Context, projectId uuid.UUID, databaseId uuid.UUID, entityId uuid.UUID) (models.DatabaseEntity, error) {
	slog.Info("Getting database entity", "entityId", entityId, "databaseId", databaseId, "projectId", projectId)

	entity, err := s.db.GetDatabaseEntityByProjectIdDatabaseIdAndId(c, sqlc.GetDatabaseEntityByProjectIdDatabaseIdAndIdParams{
		ID:         entityId,
		ProjectID:  projectId,
		DatabaseID: databaseId,
	})

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Info("Database entity not found", "entityId", entityId, "databaseId", databaseId, "projectId", projectId)
		return models.DatabaseEntity{}, &models.ErrorResponse{
			Message: errormsg.ErrDatabaseEntityNotFound,
			Errors:  nil,
			Status:  http.StatusBadRequest,
		}
	}

	if err != nil {
		slog.Error("Unable to get database entity", "entityId", entityId, "databaseId", databaseId, "projectId", projectId, "err", err)
		return models.DatabaseEntity{}, err
	}

	return mapper.ToDomainDatabaseEntity(entity), nil
}
