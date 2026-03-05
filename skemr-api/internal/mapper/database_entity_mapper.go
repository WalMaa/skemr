package mapper

import (
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-common/models"
)

func ToDomainDatabaseEntity(e sqlc.DatabaseEntity) models.DatabaseEntity {
	return models.DatabaseEntity{
		ID:          e.ID,
		Type:        models.DatabaseEntityType(e.EntityType),
		ParentId:    e.ParentID,
		Name:        e.Name,
		Status:      models.DatabaseEntityStatus(e.Status),
		Attributes:  ToMap(e.Attributes),
		DeletedAt:   TimePtr(&e.DeletedAt),
		FirstSeenAt: Time(&e.FirstSeenAt),
		CreatedAt:   Time(&e.CreatedAt),
	}
}

func ToDomainDatabaseEntities(entities []sqlc.DatabaseEntity) []models.DatabaseEntity {
	domainEntities := make([]models.DatabaseEntity, len(entities))
	for i, entity := range entities {
		domainEntities[i] = ToDomainDatabaseEntity(entity)
	}
	return domainEntities
}
