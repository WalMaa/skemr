package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/walmaa/skemr/db/sqlc"
	"github.com/walmaa/skemr/errormsg"
)

type DatabaseService struct {
	db sqlc.Querier
}

func NewDatabaseService(q sqlc.Querier) *DatabaseService {
	return &DatabaseService{db: q}
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

func (r *DatabaseService) CreateDatabase(c context.Context, args sqlc.CreateDatabaseParams) (sqlc.Database, error) {
	slog.Info("Creating database", "name", args)

	// Check if the project exists
	_, err := CheckProjectExists(c, r.db, args.ProjectID)
	if err != nil {
		return sqlc.Database{}, err
	}

	// Check a database with the given name already exists
	exists, err := r.db.GetDatabaseByNameAndProject(c, sqlc.GetDatabaseByNameAndProjectParams{
		ProjectID:   args.ProjectID,
		DisplayName: args.DisplayName,
	})

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("Error checking for existing database", "name", args.DisplayName, "err", err)
		return sqlc.Database{}, err
	}

	if exists != (sqlc.Database{}) {
		slog.Warn("Database already exists", "name", args.DisplayName, "project_id", args.ProjectID)
		return sqlc.Database{}, errormsg.ErrDatabaseAlreadyExists
	}

	return r.db.CreateDatabase(c, args)
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
