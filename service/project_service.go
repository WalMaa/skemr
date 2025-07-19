package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"skemr/db/sqlc"
)

type ProjectService struct {
	db sqlc.Querier
}

func NewProjectService(q sqlc.Querier) *ProjectService {
	return &ProjectService{db: q}
}

func (r *ProjectService) CreateProject(c *gin.Context, name string) (sqlc.Project, error) {
	slog.Info("Creating project", "name", name)
	return r.db.CreateProject(c, name)
}

func (r *ProjectService) GetProject(c *gin.Context, id uuid.UUID) (sqlc.Project, error) {
	return r.db.GetProject(c, &id)
}

func (r *ProjectService) DeleteProject(c *gin.Context, id uuid.UUID) error {
	return r.db.DeleteProject(c, &id)
}
