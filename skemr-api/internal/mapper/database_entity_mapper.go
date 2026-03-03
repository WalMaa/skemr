package mapper

import (
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-common/models"
)

func ToDomainDatabaseEntity(e sqlc.DatabaseEntity) models.DatabaseEntity {
	return models.DatabaseEntity{
		ID:         e.ID,
		ProjectId:  e.ProjectID,
		Type:       models.DatabaseEntityType(e.EntityType),
		ParentId:   e.ParentID,
		Name:       e.Name,
		Attributes: ToMap(e.Attributes),
		CreatedAt:  Time(&e.CreatedAt),
	}
}

func ToDomainDatabaseEntities(entities []sqlc.DatabaseEntity) []models.DatabaseEntity {
	domainEntities := make([]models.DatabaseEntity, len(entities))
	for i, entity := range entities {
		domainEntities[i] = ToDomainDatabaseEntity(entity)
	}
	return domainEntities
}
