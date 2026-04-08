package dbreflect

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
	connector := NewPostgresConnector(dbModel)

	mockDB := mocks.NewMockQuerier(t)

	mockDB.On("GetDatabaseEntityByDatabaseIdAndTypeAndParentAndName", mock.Anything, mock.Anything).Return(sqlc.DatabaseEntity{
		ID:         uuid.New(),
		Name:       "testSchemaName",
		DatabaseID: uuid.UUID{},
	}, nil)
	mockDB.On("GetDatabaseEntitiesByDatabaseId", mock.Anything, mock.Anything).Return([]sqlc.DatabaseEntity{
		{
			ID:         uuid.New(),
			Name:       "testSchemaName",
			DatabaseID: dbID,
		},
	}, nil)
	mockDB.On("UpdateDatabaseEntityAsDeleted", mock.Anything, mock.Anything).Return(nil)
	syncService := NewSchemaSyncService(mockDB, func(_ models.Database) DatabaseConnector { return connector })

	err = syncService.SyncSchema(ctx, dbModel)
	require.NoError(t, err)

	mockDB.AssertExpectations(t)

}

func TestUpdateSchemaCreatesNew(t *testing.T) {
	c := context.Background()
	dataBaseId := uuid.New()
	schemaName := "testSchemaName"
	mockConnector := new(MockPostgresConnector)

	database := models.Database{
		ID:        dataBaseId,
		ProjectID: uuid.New(),
	}

	schemaRef := SchemaRef{
		Name: schemaName,
	}

	mockDB := mocks.NewMockQuerier(t)
	mockDB.On("GetDatabaseEntityByDatabaseIdAndTypeAndParentAndName", mock.Anything, mock.Anything).Return(sqlc.DatabaseEntity{}, pgx.ErrNoRows)
	mockDB.On("GetDatabaseEntityByFingerprint", mock.Anything, mock.Anything).Return(sqlc.DatabaseEntity{}, pgx.ErrNoRows)
	mockDB.On("CreateDatabaseEntity", mock.Anything, mock.Anything).Return(sqlc.DatabaseEntity{
		ID:         uuid.New(),
		Name:       schemaName,
		DatabaseID: dataBaseId,
		Fingerprint: pgtype.Text{
			String: "fingerprint",
			Valid:  true,
		},
	}, nil)

	syncService := NewSchemaSyncService(mockDB, func(_ models.Database) DatabaseConnector { return mockConnector })
	schema, err := syncService.updateSchema(c, schemaRef, database)

	require.NoError(t, err)
	require.Equal(t, schema.Name, schemaName)
	mockDB.AssertExpectations(t)

}

func TestUpdateSchemaUpdatesExisting(t *testing.T) {
	c := context.Background()
	dataBaseId := uuid.New()
	schemaName := "testSchemaName"
	mockConnector := new(MockPostgresConnector)

	database := models.Database{
		ID:        dataBaseId,
		ProjectID: uuid.New(),
	}

	schemaRef := SchemaRef{
		Name: schemaName,
	}

	mockDB := mocks.NewMockQuerier(t)
	mockDB.On("GetDatabaseEntityByDatabaseIdAndTypeAndParentAndName", mock.Anything, mock.Anything).Return(sqlc.DatabaseEntity{
		ID:         uuid.New(),
		Name:       schemaName,
		DatabaseID: dataBaseId,
	}, nil)

	syncService := NewSchemaSyncService(mockDB, func(_ models.Database) DatabaseConnector { return mockConnector })
	schema, err := syncService.updateSchema(c, schemaRef, database)

	require.NoError(t, err)
	require.Equal(t, schema.Name, schemaName)
	mockDB.AssertExpectations(t)
}

func TestMarkEntityAsDeletedWhenEntityNotAppearingInSync(t *testing.T) {
	c := context.Background()
	entityId := uuid.New()
	entityId2 := uuid.New()
	databaseId := uuid.New()
	projectId := uuid.New()

	mockConnector := new(MockPostgresConnector)

	oldEntity1 := sqlc.DatabaseEntity{
		ID:         entityId,
		Name:       "testSchemaName",
		DatabaseID: databaseId,
	}
	oldEntity2 := sqlc.DatabaseEntity{
		ID:         entityId2,
		Name:       "testSchemaName2",
		DatabaseID: databaseId,
	}

	oldEntities := []sqlc.DatabaseEntity{
		oldEntity1,
		oldEntity2,
	}

	mockDB := mocks.NewMockQuerier(t)
	mockDB.On("GetDatabaseEntitiesByDatabaseId", mock.Anything, mock.Anything).Return(oldEntities, nil)
	mockDB.On("UpdateDatabaseEntityAsDeleted", mock.Anything, entityId2).Return(nil).Once()
	mockDB.On("GetDatabaseEntityByDatabaseIdAndTypeAndParentAndName", mock.Anything, mock.Anything).Return(oldEntity1, nil)

	mockConnector.On("GetTablesInSchema", mock.Anything, mock.Anything, mock.Anything).Return([]TableRef{}, nil)
	// Return only the first entity, simulating that the second one has been removed in the database and should be marked as deleted
	mockConnector.On("GetSchemas", mock.Anything, mock.Anything).Return([]SchemaRef{
		{
			Name: "testSchemaName",
		},
	}, nil)
	mockConnector.On("Connect", mock.Anything).Return(nil, nil)
	mockConnector.On("Disconnect", mock.Anything, mock.Anything).Return(nil)

	syncService := NewSchemaSyncService(mockDB, func(_ models.Database) DatabaseConnector { return mockConnector })
	err := syncService.SyncSchema(c, models.Database{ID: databaseId, ProjectID: projectId})

	require.NoError(t, err)
	mockDB.AssertCalled(t, "UpdateDatabaseEntityAsDeleted", c, entityId2)
	mockDB.AssertExpectations(t)
}
