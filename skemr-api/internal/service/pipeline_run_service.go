package service

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/dto"
	"github.com/walmaa/skemr-api/internal/errormsg"
	"github.com/walmaa/skemr-api/internal/mapper"
	"github.com/walmaa/skemr-common/models"
)

type PipelineRunStore interface {
	GetPipelineRunByDatabaseIdAndId(ctx context.Context, arg sqlc.GetPipelineRunByDatabaseIdAndIdParams) (sqlc.PipelineRun, error)
	GetPipelineRunsByDatabaseId(ctx context.Context, databaseId uuid.UUID) ([]sqlc.PipelineRun, error)
	CreatePipelineRun(ctx context.Context, pipelineRun sqlc.CreatePipelineRunParams) (sqlc.PipelineRun, error)
}

type PipelineRunService struct {
	store         PipelineRunStore
	scopeResolver ScopeResolver
}

func NewPipelineRunService(store PipelineRunStore, resolver ScopeResolver) *PipelineRunService {
	return &PipelineRunService{store: store, scopeResolver: resolver}
}

func (s *PipelineRunService) GetPipelineRun(c context.Context, projectId uuid.UUID, databaseId uuid.UUID, id uuid.UUID) (models.PipelineRun, error) {
	slog.Info("Getting pipeline run", "id", id, "databaseId", databaseId, "projectId", projectId)
	_, err := s.scopeResolver.RequireDatabase(c, projectId, databaseId)
	if err != nil {
		slog.Error("Error fetching database", "err", err)
		return models.PipelineRun{}, err
	}

	pipelineRun, err := s.store.GetPipelineRunByDatabaseIdAndId(c, sqlc.GetPipelineRunByDatabaseIdAndIdParams{
		DatabaseID: databaseId,
		ID:         id,
	})

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("Error fetching pipeline run", "err", err)
		return models.PipelineRun{}, &models.ErrorResponse{
			Errors:  nil,
			Status:  http.StatusInternalServerError,
			Message: errormsg.ErrPipelineRunFetchFailed,
		}
	}

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Info("Pipeline run not found", "id", id, "databaseId", databaseId, "projectId", projectId)
		return models.PipelineRun{}, &models.ErrorResponse{
			Errors:  nil,
			Status:  http.StatusNotFound,
			Message: errormsg.ErrPipelineRunNotFound,
		}
	}

	return mapper.ToDomainPipelineRun(pipelineRun), nil
}

func (s *PipelineRunService) GetPipelineRuns(c context.Context, projectId uuid.UUID, databaseId uuid.UUID) ([]models.PipelineRun, error) {
	slog.Info("Getting pipeline runs", "databaseId", databaseId, "projectId", projectId)
	_, err := s.scopeResolver.RequireDatabase(c, projectId, databaseId)
	if err != nil {
		slog.Error("Error fetching database", "err", err)
		return nil, err
	}

	pipelineRuns, err := s.store.GetPipelineRunsByDatabaseId(c, databaseId)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("Error fetching pipeline runs", "err", err)
		return nil, &models.ErrorResponse{
			Errors:  nil,
			Status:  http.StatusInternalServerError,
			Message: errormsg.ErrPipelineRunFetchFailed,
		}
	}

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Info("No pipeline runs found", "databaseId", databaseId, "projectId", projectId)
		return nil, &models.ErrorResponse{
			Errors:  nil,
			Status:  http.StatusNotFound,
			Message: errormsg.ErrPipelineRunNotFound,
		}
	}

	return mapper.ToDomainPipelineRuns(pipelineRuns), nil
}

func (s *PipelineRunService) CreatePipelineRun(c context.Context, projectId uuid.UUID, databaseId uuid.UUID, pipelineRun dto.PipelineRunCreationDto) (sqlc.PipelineRun, error) {
	slog.Info("Creating pipeline run", "pipelineRun", pipelineRun, "databaseId", databaseId, "projectId", projectId)
	_, err := s.scopeResolver.RequireDatabase(c, projectId, databaseId)
	if err != nil {
		slog.Error("Error fetching database", "err", err)
		return sqlc.PipelineRun{}, err
	}

	completedAt, err := time.Parse(time.RFC3339, pipelineRun.CompletedAt)

	if err != nil {
		slog.Error("Error parsing completedAt", "err", err)
		return sqlc.PipelineRun{}, err
	}

	createdPipelineRun, err := s.store.CreatePipelineRun(c, sqlc.CreatePipelineRunParams{
		DatabaseID: databaseId,
		Status:     sqlc.MigrationStatus(pipelineRun.Status),
		Environment: pgtype.Text{
			Valid:  true,
			String: pipelineRun.Environment,
		},
		CompletedAt: pgtype.Timestamptz{
			Valid: true,
			Time:  completedAt,
		},
	})

	if err != nil {
		slog.Error("Error creating pipeline run", "err", err)
		return sqlc.PipelineRun{}, err
	}

	return createdPipelineRun, nil
}
