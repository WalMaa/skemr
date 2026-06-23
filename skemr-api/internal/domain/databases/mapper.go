package databases

import (
	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/mapper"
	"github.com/walmaa/skemr-common/models"
)

func ToDomainDatabase(e sqlc.Database) models.Database {
	return models.Database{
		ID:                       e.ID,
		DisplayName:              e.DisplayName,
		DbName:                   &e.DbName.String,
		Username:                 &e.Username.String,
		Password:                 &e.Password.String,
		Host:                     &e.Host.String,
		Port:                     e.Port.Int32,
		SslMode:                  e.SslMode,
		DatabaseType:             models.DatabaseType(e.DatabaseType.DatabaseType),
		ProjectID:                e.ProjectID,
		FailedConnectionAttempts: e.FailedConnectionAttempts,
		LastSyncedAt:             mapper.TimePtr(&e.LastSyncedAt),
		LastSyncError:            mapper.TextPtr(&e.LastSyncError),
	}
}

func ToDomainDatabases(d []sqlc.Database) []models.Database {
	databases := make([]models.Database, len(d))
	for i, database := range d {
		databases[i] = ToDomainDatabase(database)
	}
	return databases
}

func ToUpdateDatabaseParams(databaseId uuid.UUID, dto DatabaseUpdateDto) sqlc.UpdateDatabaseParams {
	return sqlc.UpdateDatabaseParams{
		DatabaseID:  databaseId,
		DisplayName: mapper.Text(dto.DisplayName),
		DbName:      mapper.Text(dto.DbName),
		Username:    mapper.Text(dto.Username),
		Password:    mapper.Text(dto.Password),
		Host:        mapper.Text(dto.Host),
		Port:        mapper.Int4(dto.Port),
		SslMode:     mapper.Text(dto.SslMode),
	}
}

func ToCreateDatabaseParams(projectId uuid.UUID, dto DatabaseCreationDto) sqlc.CreateDatabaseParams {
	return sqlc.CreateDatabaseParams{
		ProjectID:    projectId,
		DisplayName:  dto.DisplayName,
		DbName:       mapper.Text(dto.DbName),
		Username:     mapper.Text(dto.Username),
		Password:     mapper.Text(dto.Password),
		Host:         mapper.Text(dto.Host),
		Port:         mapper.Int4(&dto.Port),
		DatabaseType: mapper.NullDatabaseType(dto.DatabaseType),
	}
}
