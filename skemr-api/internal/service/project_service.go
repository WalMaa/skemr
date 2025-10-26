package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/errormsg"
)

type ProjectService struct {
	db sqlc.Querier
}

func NewProjectService(q sqlc.Querier) *ProjectService {
	return &ProjectService{db: q}
}

// CheckProjectExists checks if a project with the given ID exists in the database.
// Used when validating the operation on resources that are tied to a project.
func CheckProjectExists(c context.Context, db sqlc.Querier, projectID uuid.UUID) (sqlc.Project, error) {
	slog.Info("Checking if project exists", "project_id", projectID)

	// Check if the project exists
	project, err := db.GetProject(c, projectID)
	if err != nil {
		slog.Error("Error getting project", "project_id", projectID, "err", err)
		return sqlc.Project{}, errormsg.ErrProjectNotFound
	}

	return project, nil
}

func (r *ProjectService) CreateProject(c context.Context, name string) (sqlc.Project, error) {

	slog.Info("Creating project", "name", name)
	return r.db.CreateProject(c, name)
}

func (r *ProjectService) GetProjects(c context.Context) ([]sqlc.Project, error) {
	slog.Info("Fetching all projects")
	projects, err := r.db.GetProjects(c)
	if err != nil {
		slog.Error("Error fetching projects", "err", err)
		return nil, err
	}
	return projects, nil
}

func (r *ProjectService) GetProject(c context.Context, id uuid.UUID) (sqlc.Project, error) {
	return r.db.GetProject(c, id)
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
