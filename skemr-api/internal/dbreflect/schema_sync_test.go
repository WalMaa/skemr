package dbreflect

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/test/mocks"
	"github.com/walmaa/skemr-common/models"
)

func TestSchemaSync(t *testing.T) {
	ctx := context.Background()
	pgC := GetTestPostgres(t)

	host, err := pgC.Host(ctx)
	require.NoError(t, err)

	port, err := pgC.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dbID := uuid.New()
	hostStr := host
	dbName := "postgres"
	username := "user"
	password := "password"

	dbModel := models.Database{
		ID:           dbID,
		DisplayName:  "Test Database",
		Username:     &username,
		Password:     &password,
		Host:         &hostStr,
		DbName:       &dbName,
		Port:         int32(port.Int()),
		DatabaseType: models.Postgres,
		ProjectID:    uuid.New(),
	}

	mockDB := mocks.NewMockQuerier(t)

	mockDB.On("GetDatabaseEntityByDatabaseIdAndTypeAndName", mock.Anything, mock.Anything).Return(sqlc.DatabaseEntity{
		ID:         uuid.New(),
		Name:       "testSchemaName",
		DatabaseID: uuid.UUID{},
	}, nil)
	syncService := NewSchemaSyncService(mockDB)

	err = syncService.SyncSchema(ctx, dbModel)
	require.NoError(t, err)

	mockDB.AssertExpectations(t)

}

func TestUpdateSchemaCreatesNew(t *testing.T) {
	c := context.Background()
	dataBaseId := uuid.New()
	schemaName := "testSchemaName"

	database := models.Database{
		ID:        dataBaseId,
		ProjectID: uuid.New(),
	}

	mockDB := mocks.NewMockQuerier(t)
	mockDB.On("GetDatabaseEntityByDatabaseIdAndTypeAndName", mock.Anything, mock.Anything).Return(sqlc.DatabaseEntity{}, pgx.ErrNoRows)
	mockDB.On("CreateDatabaseEntity", mock.Anything, mock.Anything).Return(sqlc.DatabaseEntity{
		ID:         uuid.New(),
		Name:       schemaName,
		DatabaseID: dataBaseId,
	}, nil)

	syncService := NewSchemaSyncService(mockDB)
	schema, err := syncService.updateSchema(c, schemaName, database)

	require.NoError(t, err)
	require.Equal(t, schema.Name, schemaName)
	mockDB.AssertExpectations(t)

}

func TestUpdateSchemaUpdatesExisting(t *testing.T) {
	c := context.Background()
	dataBaseId := uuid.New()
	schemaName := "testSchemaName"

	database := models.Database{
		ID:        dataBaseId,
		ProjectID: uuid.New(),
	}

	mockDB := mocks.NewMockQuerier(t)
	mockDB.On("GetDatabaseEntityByDatabaseIdAndTypeAndName", mock.Anything, mock.Anything).Return(sqlc.DatabaseEntity{
		ID:         uuid.New(),
		Name:       schemaName,
		DatabaseID: dataBaseId,
	}, nil)

	syncService := NewSchemaSyncService(mockDB)
	schema, err := syncService.updateSchema(c, schemaName, database)

	require.NoError(t, err)
	require.Equal(t, schema.Name, schemaName)
	mockDB.AssertExpectations(t)

}
