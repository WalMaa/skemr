package service

import (
	"github.com/gin-gonic/gin"
	skemr "skemr/db"
)

type ProjectService struct {
	Queries *skemr.Queries
}

func NewProjectService(q *skemr.Queries) *ProjectService {
	return &ProjectService{Queries: q}
}

func (r *ProjectService) CreateProject(c *gin.Context, name string) (skemr.Project, error) {
	return r.Queries.CreateProject(c, name)

}
