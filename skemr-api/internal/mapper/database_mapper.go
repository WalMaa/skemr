package mapper

import (
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-common/models"
)

func ToDomainDatabase(e sqlc.Database) models.Database {
	return models.Database{
		ID:           e.ID,
		DisplayName:  e.DisplayName,
		DbName:       &e.DbName.String,
		Username:     &e.Username.String,
		Password:     &e.Password.String,
		Host:         &e.Host.String,
		Port:         e.Port.Int32,
		DatabaseType: models.DatabaseType(e.DatabaseType.DatabaseType),
		ProjectID:    e.ProjectID,
	}
}
