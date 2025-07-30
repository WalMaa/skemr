package service

import (
	"errors"
	"github.com/gin-gonic/gin"
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

func checkProjectExists(c *gin.Context, db sqlc.Querier, projectID uuid.UUID) (sqlc.Project, error) {
	slog.Info("Checking if project exists", "project_id", projectID)

	// Check if the project exists
	project, err := db.GetProject(c, projectID)
	if err != nil {
		slog.Error("Error getting project", "project_id", projectID, "err", err)
		return sqlc.Project{}, errormsg.ErrProjectNotFound
	}

	return project, nil
}

func (r *DatabaseService) CreateDatabase(c *gin.Context, args sqlc.CreateDatabaseParams) (sqlc.Database, error) {
	slog.Info("Creating database", "name", args)

	// Check if the project exists
	_, err := checkProjectExists(c, r.db, args.ProjectID)
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

func (r *DatabaseService) GetDatabase(c *gin.Context, id uuid.UUID) (sqlc.Database, error) {
	slog.Info("Getting database", "id", id)
	return r.db.GetDatabase(c, id)
}

func (r *DatabaseService) DeleteDatabase(c *gin.Context, id uuid.UUID) error {
	slog.Info("Deleting database", "id", id)
	return r.db.DeleteDatabase(c, id)
}

func (r *DatabaseService) ListDatabasesByProject(c *gin.Context, id uuid.UUID) ([]sqlc.Database, error) {
	slog.Info("Listing databases for project", "project_id", id)
	return r.db.ListDatabasesByProject(c, id)
}
