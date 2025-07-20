package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"skemr/db/sqlc"
)

type DatabaseService struct {
	db sqlc.Querier
}

func NewDatabaseService(q sqlc.Querier) *DatabaseService {
	return &DatabaseService{db: q}
}

func (r *DatabaseService) CreateDatabase(c *gin.Context, args sqlc.CreateDatabaseParams) (sqlc.Database, error) {
	slog.Info("Creating database", "name", args)
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
