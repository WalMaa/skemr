package mapper

import (
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-common/models"
)

func ToDomainDatabaseChange(change sqlc.DatabaseChange) models.DatabaseChange {
	return models.DatabaseChange{
		Id:         change.ID,
		DatabaseId: change.DatabaseID,
		EntityId:   change.EntityID,
		Action:     models.MigrationStatementAction(change.Action),
		CreatedAt:  Time(&change.CreatedAt),
	}
}

func ToDomainDatabaseChanges(changes []sqlc.DatabaseChange) []models.DatabaseChange {
	databaseChanges := make([]models.DatabaseChange, len(changes))
	for i, change := range changes {
		databaseChanges[i] = ToDomainDatabaseChange(change)
	}
	return databaseChanges
}
