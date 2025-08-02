package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/walmaa/skemr/db/sqlc"
	"github.com/walmaa/skemr/errormsg"
	"log/slog"
)

type DatabaseService struct {
	db sqlc.Querier
}

func NewDatabaseService(q sqlc.Querier) *DatabaseService {
	return &DatabaseService{db: q}
}

func (r *DatabaseService) CreateDatabase(c context.Context, args sqlc.CreateDatabaseParams) (sqlc.Database, error) {
	slog.Info("Creating database", "name", args)

	// Check if the project exists
	_, err := CheckProjectExists(c, r.db, args.ProjectID)
	if err != nil {
		return sqlc.Database{}, err
	}

	// Check a database with the given name already exists
	exists, err := r.db.GetDatabaseByName(c, sqlc.GetDatabaseByNameParams{
		ProjectID: args.ProjectID,
		Name:      args.Name,
	})

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("Error checking for existing database", "name", args.Name, "err", err)
		return sqlc.Database{}, err
	}

	if exists != (sqlc.Database{}) {
		slog.Warn("Database already exists", "name", args.Name, "project_id", args.ProjectID)
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
