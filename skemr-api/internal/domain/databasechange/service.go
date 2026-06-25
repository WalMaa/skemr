package databasechange

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/errormsg"
	"github.com/walmaa/skemr-api/internal/service"
	"github.com/walmaa/skemr-common/models"
)

type DatabaseChangeStore interface {
	GetDatabaseChangeByDatabaseIdAndId(ctx context.Context, arg sqlc.GetDatabaseChangeByDatabaseIdAndIdParams) (sqlc.DatabaseChange, error)
	GetDatabaseChangesByDatabaseIdAndId(c context.Context, params sqlc.GetDatabaseChangesByDatabaseIdAndIdParams) ([]sqlc.DatabaseChange, error)
}

type DatabaseChangeService struct {
	store         DatabaseChangeStore
	scopeResolver service.ScopeResolver
}

func NewDatabaseChangeService(store DatabaseChangeStore, resolver service.ScopeResolver) *DatabaseChangeService {
	return &DatabaseChangeService{store: store, scopeResolver: resolver}
}

func (s *DatabaseChangeService) GetDatabaseChange(c context.Context, projectId uuid.UUID, databaseId uuid.UUID, id uuid.UUID) (models.DatabaseChange, error) {
	slog.Info("Getting database change", "id", id, "projectId", projectId, "databaseId", databaseId)

	_, err := s.scopeResolver.RequireDatabase(c, projectId, databaseId)
	if err != nil {
		slog.Error("Error fetching database", "err", err)
		return models.DatabaseChange{}, err
	}

	databaseChange, err := s.store.GetDatabaseChangeByDatabaseIdAndId(c, sqlc.GetDatabaseChangeByDatabaseIdAndIdParams{
		ID:         id,
		DatabaseID: databaseId,
	})

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("Error fetching database change", "err", err)
		return models.DatabaseChange{}, &models.ErrorResponse{}
	}

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Info("Database change not found", "id", id, "projectId", projectId, "databaseId", databaseId)
		return models.DatabaseChange{}, &models.ErrorResponse{
			Message: errormsg.ErrDatabaseChangeNotFound,
			Status:  http.StatusNotFound,
			Errors:  nil,
		}
	}

	return ToDomainDatabaseChange(databaseChange), nil
}

func (s *DatabaseChangeService) GetDatabaseChanges(c context.Context, projectId uuid.UUID, databaseId uuid.UUID, limit int, offset int) ([]models.DatabaseChange, error) {
	slog.Info("Getting database changes", "limit", limit, "offset", offset, "projectId", projectId, "databaseId", databaseId)

	_, err := s.scopeResolver.RequireDatabase(c, projectId, databaseId)
	if err != nil {
		slog.Error("Error fetching database", "err", err)
		return nil, err
	}

	databaseChanges, err := s.store.GetDatabaseChangesByDatabaseIdAndId(c, sqlc.GetDatabaseChangesByDatabaseIdAndIdParams{
		DatabaseID: databaseId,
		Offset: pgtype.Int4{
			Int32: int32(offset),
			Valid: true,
		},
		Limit: pgtype.Int4{
			Int32: int32(limit),
			Valid: true,
		},
	})

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("Error fetching database changes", "err", err)
		return nil, &models.ErrorResponse{
			Errors:  nil,
			Status:  http.StatusInternalServerError,
			Message: "Error fetching database changes",
		}
	}

	return ToDomainDatabaseChanges(databaseChanges), nil
}
