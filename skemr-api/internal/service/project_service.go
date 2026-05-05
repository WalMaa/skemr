package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/dto"
	"github.com/walmaa/skemr-api/internal/errormsg"
	"github.com/walmaa/skemr-api/internal/mapper"
	"github.com/walmaa/skemr-common/models"
)

type ProjectService struct {
	db sqlc.Querier
}

func NewProjectService(q sqlc.Querier) *ProjectService {
	return &ProjectService{db: q}
}

// CheckProjectExists checks if a project with the given ID exists in the database.
// Used when validating the operation on resources that are tied to a project.
func CheckProjectExists(c context.Context, db sqlc.Querier, projectID uuid.UUID) (models.Project, error) {
	slog.Info("Checking if project exists", "project_id", projectID)

	// Check if the project exists
	project, err := db.GetProject(c, projectID)
	if err != nil {
		slog.Error("Error getting project", "project_id", projectID, "err", err)
		return models.Project{}, &models.ErrorResponse{
			Message: errormsg.ErrProjectNotFound,
			Errors:  nil,
			Status:  404,
		}
	}

	return mapper.ToDomainProject(project), nil
}

func (r *ProjectService) CreateProject(c context.Context, dto dto.ProjectCreationDto) (models.Project, error) {

	slog.Info("Creating project", "name", dto.Name)
	project, err := r.db.CreateProject(c, dto.Name)
	if err != nil {
		slog.Error("Error creating project", "name", dto.Name, "err", err)
		return models.Project{}, err
	}
	return mapper.ToDomainProject(project), nil
}

func (r *ProjectService) GetProjects(c context.Context) ([]models.Project, error) {
	slog.Info("Fetching all projects")
	projects, err := r.db.GetProjects(c)
	if err != nil {
		slog.Error("Error fetching projects", "err", err)
		return nil, err
	}
	return mapper.ToDomainProjects(projects), nil
}

func (r *ProjectService) GetProject(c context.Context, projectId uuid.UUID) (models.Project, error) {
	slog.Info("Getting project", "projectId", projectId)
	project, err := r.db.GetProject(c, projectId)

	if err != nil {
		slog.Error("Error getting project", "err", err)
		return models.Project{}, err
	}

	return mapper.ToDomainProject(project), nil
}

func (r *ProjectService) DeleteProject(c context.Context, id uuid.UUID) error {
	slog.Info("Deleting project", "id", id)
	// Check if the project exists
	_, err := CheckProjectExists(c, r.db, id)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("Error checking if project exists", "id", id, "err", err)
		return err
	}

	return r.db.DeleteProject(c, id)
}
