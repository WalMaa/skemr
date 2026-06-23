package databasechange

import (
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/mapper"
	"github.com/walmaa/skemr-common/models"
)

func ToDomainDatabaseChange(change sqlc.DatabaseChange) models.DatabaseChange {
	return models.DatabaseChange{
		Id:        change.ID,
		EntityId:  change.EntityID,
		Action:    models.MigrationStatementAction(change.Action),
		CreatedAt: mapper.Time(&change.CreatedAt),
	}
}

func ToDomainDatabaseChanges(changes []sqlc.DatabaseChange) []models.DatabaseChange {
	databaseChanges := make([]models.DatabaseChange, len(changes))
	for i, change := range changes {
		databaseChanges[i] = ToDomainDatabaseChange(change)
	}
	return databaseChanges
}
