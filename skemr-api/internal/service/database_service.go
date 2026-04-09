package service

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/dto"
	"github.com/walmaa/skemr-api/internal/errormsg"
	"github.com/walmaa/skemr-api/internal/mapper"
	"github.com/walmaa/skemr-api/internal/tasks"
	"github.com/walmaa/skemr-common/models"
)

type DatabaseService struct {
	db         sqlc.Querier
	taskClient *asynq.Client
}

func NewDatabaseService(q sqlc.Querier, c *asynq.Client) *DatabaseService {
	return &DatabaseService{db: q, taskClient: c}
}

func CheckDatabaseExists(c context.Context, db sqlc.Querier, projectId uuid.UUID, dbId uuid.UUID) (models.Database, error) {
	slog.Info("Checking if database exists", "database_id", dbId)

	// Check if the database exists
	database, err := db.GetDatabaseByIdAndProject(c, sqlc.GetDatabaseByIdAndProjectParams{
		ID:        dbId,
		ProjectID: projectId,
	})
	if err != nil {
		slog.Error("Error getting database", "database_id", dbId, "err", err)
		return models.Database{}, &models.ErrorResponse{
			Message: errormsg.ErrDatabaseNotFound,
			Errors:  nil,
			Status:  http.StatusNotFound,
		}
	}

	return mapper.ToDomainDatabase(database), nil
}

func (r *DatabaseService) CreateDatabase(c context.Context, projectId uuid.UUID, dto dto.DatabaseCreationDto) (models.Database, error) {
	slog.Info("Creating database", "name", dto)

	// Check if the project exists
	_, err := CheckProjectExists(c, r.db, projectId)
	if err != nil {
		return models.Database{}, err
	}

	// Check a database with the given name already exists
	exists, err := r.db.GetDatabaseByNameAndProject(c, sqlc.GetDatabaseByNameAndProjectParams{
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
	database, err := r.db.CreateDatabase(c, mapper.ToCreateDatabaseParams(projectId, dto))

	if err != nil {
		slog.Error("Error creating database", err)
		return models.Database{}, err
	}

	r.createDatabaseSyncTask(database.ID)
	return mapper.ToDomainDatabase(database), nil
}

// CreateDatabaseSyncTask Creates a Database sync task for Asynq background processing.
func (r *DatabaseService) createDatabaseSyncTask(databaseId uuid.UUID) {
	// TODO: rate limiting
	slog.Info("Creating a Datbase sync task", "databaseId", databaseId)
	task, err := tasks.NewDatabaseSyncTask(databaseId)
	if err != nil {
		slog.Error("Unable to create database sync task")
	}
	_, err = r.taskClient.Enqueue(task)
	if err != nil {
		slog.Error("Error in task", err)
	}

}

func (r *DatabaseService) EnqueueManualDatabaseSync(c context.Context, projectId uuid.UUID, databaseId uuid.UUID) error {
	slog.Info("Enqueuing manual database sync", "projectId", projectId, "databaseId", databaseId)

	project, err := CheckProjectExists(c, r.db, projectId)

	if err != nil {
		slog.Error("Error fetching project")
		return err
	}

	database, err := CheckDatabaseExists(c, r.db, project.ID, databaseId)

	if err != nil {
		slog.Error("Error getting database")
		return err
	}

	r.createDatabaseSyncTask(database.ID)

	return nil
}

func (r *DatabaseService) GetDatabase(c context.Context, databaseId uuid.UUID) (models.Database, error) {
	slog.Info("Getting database", "databaseId", databaseId)
	database, err := r.db.GetDatabase(c, databaseId)
	if err != nil {
		slog.Error("Unable to get database")
		return models.Database{}, err
	}
	return mapper.ToDomainDatabase(database), nil
}

func (r *DatabaseService) DeleteDatabase(c context.Context, id uuid.UUID) error {
	slog.Info("Deleting database", "id", id)
	return r.db.DeleteDatabase(c, id)
}

func (r *DatabaseService) ListDatabasesByProject(c context.Context, projectId uuid.UUID) ([]models.Database, error) {
	slog.Info("Listing databases for project", "project_id", projectId)
	project, err := CheckProjectExists(c, r.db, projectId)
	if err != nil {
		slog.Error("Could not get project")
		return nil, err
	}
	databases, err := r.db.ListDatabasesByProject(c, project.ID)

	if err != nil {
		slog.Error("Unable to get databases")
		return nil, err
	}

	return mapper.ToDomainDatabases(databases), nil
}

func (r *DatabaseService) UpdateDatabase(c context.Context, projectId uuid.UUID, databaseId uuid.UUID, dto dto.DatabaseUpdateDto) (models.Database, error) {

	slog.Info("Updating database", "id", databaseId)

	project, err := CheckProjectExists(c, r.db, projectId)

	if err != nil {
		slog.Error("Error fetching project")
		return models.Database{}, err
	}

	_, err = CheckDatabaseExists(c, r.db, project.ID, databaseId)

	if err != nil {
		slog.Error("Error getting database")
		return models.Database{}, err
	}

	database, err := r.db.UpdateDatabase(c, mapper.ToUpdateDatabaseParams(databaseId, dto))

	if err != nil {
		slog.Error("Error updating database", err)
		return models.Database{}, err
	}

	return mapper.ToDomainDatabase(database), nil
}
