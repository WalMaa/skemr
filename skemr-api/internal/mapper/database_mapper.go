package mapper

import (
	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/dto"
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

func ToDomainDatabases(d []sqlc.Database) []models.Database {
	databases := make([]models.Database, len(d))
	for i, database := range d {
		databases[i] = ToDomainDatabase(database)
	}
	return databases
}

func ToUpdateDatabaseParams(databaseId uuid.UUID, dto dto.DatabaseUpdateDto) sqlc.UpdateDatabaseParams {
	return sqlc.UpdateDatabaseParams{
		DatabaseID:  databaseId,
		DisplayName: Text(dto.DisplayName),
		DbName:      Text(dto.DbName),
		Username:    Text(dto.Username),
		Password:    Text(dto.Password),
		Host:        Text(dto.Host),
		Port:        Int4(dto.Port),
	}
}

func ToCreateDatabaseParams(projectId uuid.UUID, dto dto.DatabaseCreationDto) sqlc.CreateDatabaseParams {
	return sqlc.CreateDatabaseParams{
		ProjectID:    projectId,
		DisplayName:  dto.DisplayName,
		DbName:       Text(dto.DbName),
		Username:     Text(dto.Username),
		Password:     Text(dto.Password),
		Host:         Text(dto.Host),
		Port:         Int4(&dto.Port),
		DatabaseType: NullDatabaseType(dto.DatabaseType),
	}
}
