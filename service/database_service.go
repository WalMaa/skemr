package service

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"log/slog"
	skemr "skemr/db"
)

type DatabaseService struct {
	Queries *skemr.Queries
}

func NewDatabaseService(q *skemr.Queries) *DatabaseService {
	return &DatabaseService{Queries: q}
}

func (r *DatabaseService) CreateDatabase(c *gin.Context, args skemr.CreateDatabaseParams) (skemr.Database, error) {
	slog.Info("Creating database", "name", args)
	return r.Queries.CreateDatabase(c, args)
}

func (r *DatabaseService) GetDatabase(c *gin.Context, id pgtype.UUID) (skemr.Database, error) {
	slog.Info("Getting database", "id", id)
	return r.Queries.GetDatabase(c, id)
}

func (r *DatabaseService) DeleteDatabase(c *gin.Context, id pgtype.UUID) error {
	slog.Info("Deleting database", "id", id)
	return r.Queries.DeleteDatabase(c, id)
}

func (r *DatabaseService) ListDatabasesByProject(c *gin.Context, id pgtype.UUID) ([]skemr.Database, error) {
	slog.Info("Listing databases for project", "project_id", id)
	return r.Queries.ListDatabasesByProject(c, id)
}
