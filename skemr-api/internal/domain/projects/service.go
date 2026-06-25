package projects

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/service"
	"github.com/walmaa/skemr-common/models"
)

type ProjectService struct {
	db            sqlc.Querier
	scopeResolver service.ScopeResolver
}

func NewProjectService(q sqlc.Querier, scopeResolver service.ScopeResolver) *ProjectService {
	return &ProjectService{db: q, scopeResolver: scopeResolver}
}

func (s *ProjectService) CreateProject(c context.Context, dto ProjectCreationDto) (models.Project, error) {

	slog.Info("Creating project", "name", dto.Name)
	project, err := s.db.CreateProject(c, dto.Name)
	if err != nil {
		slog.Error("Error creating project", "name", dto.Name, "err", err)
		return models.Project{}, err
	}
	return ToDomainProject(project), nil
}

func (r *ProjectService) GetProjects(c context.Context) ([]models.Project, error) {
	slog.Info("Fetching all projects")
	projects, err := r.db.GetProjects(c)
	if err != nil {
		slog.Error("Error fetching projects", "err", err)
		return nil, err
	}
	return ToDomainProjects(projects), nil
}

func (r *ProjectService) GetProject(c context.Context, projectId uuid.UUID) (models.Project, error) {
	slog.Info("Getting project", "projectId", projectId)
	project, err := r.db.GetProject(c, projectId)

	if err != nil {
		slog.Error("Error getting project", "err", err)
		return models.Project{}, err
	}

	return ToDomainProject(project), nil
}

func (s *ProjectService) DeleteProject(c context.Context, id uuid.UUID) error {
	slog.Info("Deleting project", "id", id)
	// Check if the project exists
	_, err := s.scopeResolver.RequireProject(c, id)

	if err != nil {
		slog.Error("Error getting project", "err", err)
		return err
	}

	return s.db.DeleteProject(c, id)
}
