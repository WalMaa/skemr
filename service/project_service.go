package service

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"log/slog"
	skemr "skemr/db"
)

type ProjectService struct {
	Queries *skemr.Queries
}

func NewProjectService(q *skemr.Queries) *ProjectService {
	return &ProjectService{Queries: q}
}

func (r *ProjectService) CreateProject(c *gin.Context, name string) (skemr.Project, error) {
	slog.Info("Creating project", "name", name)
	return r.Queries.CreateProject(c, name)
}

func (r *ProjectService) GetProject(c *gin.Context, id pgtype.UUID) (skemr.Project, error) {
	return r.Queries.GetProject(c, id)
}

func (r *ProjectService) DeleteProject(c *gin.Context, id pgtype.UUID) error {
	return r.Queries.DeleteProject(c, id)
}
