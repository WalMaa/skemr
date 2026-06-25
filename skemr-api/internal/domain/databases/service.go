package databases

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/errormsg"
	"github.com/walmaa/skemr-api/internal/service"
	"github.com/walmaa/skemr-api/internal/tasks"
	"github.com/walmaa/skemr-common/models"
)

type DatabaseService struct {
	db            sqlc.Querier
	taskClient    *asynq.Client
	scopeResolver service.ScopeResolver
}

func NewDatabaseService(q sqlc.Querier, c *asynq.Client, scopeResolver service.ScopeResolver) *DatabaseService {
	return &DatabaseService{db: q, taskClient: c, scopeResolver: scopeResolver}
}

func (s *DatabaseService) CreateDatabase(c context.Context, projectId uuid.UUID, dto DatabaseCreationDto) (models.Database, error) {
	slog.Info("Creating database", "name", dto)

	// Check if the project exists
	_, err := s.scopeResolver.RequireProject(c, projectId)
	if err != nil {
		return models.Database{}, err
	}

	// Check if a database with the given name already exists
	exists, err := s.db.GetDatabaseByNameAndProject(c, sqlc.GetDatabaseByNameAndProjectParams{
		ProjectID:   projectId,
		DisplayName: dto.DisplayName,
	})

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("Error checking for existing database", "name", dto.DisplayName, "err", err)
		return models.Database{}, err
	}

	if exists != (sqlc.Database{}) {
		slog.Warn("Database already exists", "name", dto.DisplayName, "project_id", projectId)
		return models.Database{}, &models.ErrorResponse{
			Message: errormsg.ErrDatabaseAlreadyExists,
			Status:  http.StatusConflict,
		}
	}
	database, err := s.db.CreateDatabase(c, ToCreateDatabaseParams(projectId, dto))

	if err != nil {
		slog.Error("Error creating database", "err", err)
		return models.Database{}, err
	}

	s.createDatabaseSyncTask(database.ID)
	return ToDomainDatabase(database), nil
}

// CreateDatabaseSyncTask Creates a Database sync task for Asynq background processing.
func (s *DatabaseService) createDatabaseSyncTask(databaseId uuid.UUID) {

	// TODO: rate limiting
	slog.Info("Creating a Datbaase sync task", "databaseId", databaseId)
	task, err := tasks.NewDatabaseSyncTask(databaseId)
	if err != nil {
		slog.Error("Unable to create database sync task")
	}
	_, err = s.taskClient.Enqueue(task)
	if err != nil {
		slog.Error("Error in task", "err", err)
	}

}

func (s *DatabaseService) EnqueueManualDatabaseSync(c context.Context, projectId uuid.UUID, databaseId uuid.UUID) error {
	slog.Info("Enqueuing manual database sync", "projectId", projectId, "databaseId", databaseId)

	database, err := s.scopeResolver.RequireDatabase(c, projectId, databaseId)

	if err != nil {
		slog.Error("Error getting database", "err", err)
		return err
	}

	s.createDatabaseSyncTask(database.ID)

	return nil
}

func (s *DatabaseService) GetDatabase(c context.Context, projectId uuid.UUID, databaseId uuid.UUID) (models.Database, error) {
	slog.Info("Getting database", "databaseId", databaseId)

	project, err := s.scopeResolver.RequireProject(c, projectId)
	if err != nil {
		slog.Error("Error fetching project", "err", err)
		return models.Database{}, err
	}

	database, err := s.db.GetDatabaseByIDAndProjectID(c, sqlc.GetDatabaseByIDAndProjectIDParams{
		ID:        databaseId,
		ProjectID: project.ID,
	})

	if err != nil {
		slog.Error("Unable to get database", "databaseId", databaseId, "err", err)
		return models.Database{}, err
	}

	return ToDomainDatabase(database), nil
}

func (s *DatabaseService) DeleteDatabase(c context.Context, id uuid.UUID) error {
	slog.Info("Deleting database", "id", id)
	return s.db.DeleteDatabase(c, id)
}

func (s *DatabaseService) ListDatabasesByProject(c context.Context, projectId uuid.UUID) ([]models.Database, error) {
	slog.Info("Listing databases for project", "project_id", projectId)
	project, err := s.scopeResolver.RequireProject(c, projectId)
	if err != nil {
		slog.Error("Could not get project")
		return nil, err
	}
	databases, err := s.db.GetDatabasesByProjectId(c, project.ID)

	if err != nil {
		slog.Error("Unable to get databases", "project_id", projectId, "err", err)
		return nil, err
	}

	return ToDomainDatabases(databases), nil
}

func (s *DatabaseService) UpdateDatabase(c context.Context, projectId uuid.UUID, databaseId uuid.UUID, dto DatabaseUpdateDto) (models.Database, error) {

	slog.Info("Updating database", "id", databaseId)

	database, err := s.scopeResolver.RequireDatabase(c, projectId, databaseId)

	if err != nil {
		slog.Error("Error getting database")
		return models.Database{}, err
	}

	updatedDatabase, err := s.db.UpdateDatabase(c, ToUpdateDatabaseParams(database.ID, dto))

	if err != nil {
		slog.Error("Error updating database", "err", err)
		return models.Database{}, err
	}

	return ToDomainDatabase(updatedDatabase), nil
}
