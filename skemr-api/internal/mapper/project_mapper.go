package mapper

import (
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-common/models"
)

func ToDomainProject(project sqlc.Project) models.Project {
	return models.Project{
		ID:        project.ID,
		Name:      project.Name,
		CreatedAt: Time(&project.CreatedAt),
		UpdatedAt: Time(&project.UpdatedAt),
	}
}

func ToDomainProjects(p []sqlc.Project) []models.Project {
	projects := make([]models.Project, len(p))
	for i, project := range p {
		projects[i] = ToDomainProject(project)
	}
	return projects
}
