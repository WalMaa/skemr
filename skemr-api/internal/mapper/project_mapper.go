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
