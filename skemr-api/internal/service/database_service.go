package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/errormsg"
	"github.com/walmaa/skemr-api/internal/dto"
	"github.com/walmaa/skemr-api/internal/mapper"
	"github.com/walmaa/skemr-api/tasks"
	"github.com/walmaa/skemr-common/models"
)

type DatabaseService struct {
	db         sqlc.Querier
	taskClient *asynq.Client
}

func NewDatabaseService(q sqlc.Querier, c *asynq.Client) *DatabaseService {
	return &DatabaseService{db: q, taskClient: c}
}

func CheckDatabaseExists(c context.Context, db sqlc.Querier, projectId uuid.UUID, dbId uuid.UUID) (sqlc.Database, error) {
	slog.Info("Checking if database exists", "database_id", dbId)

	// Check if the database exists
	database, err := db.GetDatabaseByIdAndProject(c, sqlc.GetDatabaseByIdAndProjectParams{
		ID:        dbId,
		ProjectID: projectId,
	})
	if err != nil {
		slog.Error("Error getting database", "database_id", dbId, "err", err)
		return sqlc.Database{}, errormsg.ErrDatabaseNotFound
	}

	return database, nil
}

func (r *DatabaseService) CreateDatabase(c context.Context, args sqlc.CreateDatabaseParams) (models.Database, error) {
	slog.Info("Creating database", "name", args)

	// Check if the project exists
	_, err := CheckProjectExists(c, r.db, args.ProjectID)
	if err != nil {
		return models.Database{}, err
	}

	// Check a database with the given name already exists
	exists, err := r.db.GetDatabaseByNameAndProject(c, sqlc.GetDatabaseByNameAndProjectParams{
		ProjectID:   args.ProjectID,
		DisplayName: args.DisplayName,
	})

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("Error checking for existing database", "name", args.DisplayName, "err", err)
		return models.Database{}, err
	}

	if exists != (sqlc.Database{}) {
		slog.Warn("Database already exists", "name", args.DisplayName, "project_id", args.ProjectID)
		return models.Database{}, errormsg.ErrDatabaseAlreadyExists
	}
	database, err := r.db.CreateDatabase(c, args)

	if err != nil {
		slog.Error("Error creating database", err)
		return models.Database{}, err
	}

	task, err := tasks.NewDatabaseSyncTask(database.ID)
	if err != nil {
		slog.Error("Unable to create database sync task")
	}
	_, err = r.taskClient.Enqueue(task)
	if err != nil {
		slog.Error("Error in task", err)
	}

	return mapper.ToDomainDatabase(database), nil
}

func (r *DatabaseService) GetDatabase(c context.Context, id uuid.UUID) (sqlc.Database, error) {
	slog.Info("Getting database", "id", id)
	return r.db.GetDatabase(c, id)
}

func (r *DatabaseService) DeleteDatabase(c context.Context, id uuid.UUID) error {
	slog.Info("Deleting database", "id", id)
	return r.db.DeleteDatabase(c, id)
}

func (r *DatabaseService) ListDatabasesByProject(c context.Context, id uuid.UUID) ([]sqlc.Database, error) {
	slog.Info("Listing databases for project", "project_id", id)
	return r.db.ListDatabasesByProject(c, id)
}

func (r *DatabaseService) UpdateDatabase(c context.Context, projectId uuid.UUID, databaseId uuid.UUID, dto dto.DatabaseUpdateDto) (models.Database, error) {

	slog.Info("Updating database", "id", databaseId)

	database, err := r.db.UpdateDatabase(c, mapper.ToUpdateDatabaseParams(databaseId, dto))

	if err != nil {
		slog.Error("Error updating database", err)
		return models.Database{}, err
	}

	return mapper.ToDomainDatabase(database), nil
}
