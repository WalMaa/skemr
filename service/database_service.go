package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"skemr/db/sqlc"
	errors "skemr/erros"
)

type DatabaseService struct {
	db sqlc.Querier
}

func NewDatabaseService(q sqlc.Querier) *DatabaseService {
	return &DatabaseService{db: q}
}

func (r *DatabaseService) CreateDatabase(c *gin.Context, args sqlc.CreateDatabaseParams) (sqlc.Database, error) {
	slog.Info("Creating database", "name", args)

	exists, err := r.db.GetDatabaseByName(c, sqlc.GetDatabaseByNameParams{
		ProjectID: args.ProjectID,
		Name:      args.Name,
	})

	if err != nil {
		slog.Error("Error checking for existing database", "name", args.Name, "err", err)
	}

	if exists != (sqlc.Database{}) {
		slog.Warn("Database already exists", "name", args.Name, "project_id", args.ProjectID)
		return sqlc.Database{}, errors.ErrDatabaseAlreadyExists
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
