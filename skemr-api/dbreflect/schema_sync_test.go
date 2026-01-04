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
)

func TestSchemaSync(t *testing.T) {
	ctx := context.Background()
	pgC := GetTestPostgres(t)

	host, err := pgC.Host(ctx)
	require.NoError(t, err)

	port, err := pgC.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dataBaseId := uuid.New()

	dbModel := &sqlc.Database{
		ID:          dataBaseId,
		DisplayName: "Test Database",
		Username:    pgtype.Text{String: "user", Valid: true},
		Password:    pgtype.Text{String: "password", Valid: true},
		Host:        pgtype.Text{String: host, Valid: true},
		Type:        "postgres",
		DbName:      "postgres",
		Port:        int32(port.Int()),
		ProjectID:   uuid.New(),
	}

	mockDB := mocks.NewMockQuerier(t)

	mockDB.On("GetDatabase", mock.Anything, mock.Anything).Return(*dbModel, nil)

	mockDB.On("GetDatabaseEntityByDatabaseIdAndTypeAndName", mock.Anything, mock.Anything).Return(sqlc.DatabaseEntity{
		ID:         uuid.New(),
		Name:       "testSchemaName",
		DatabaseID: uuid.UUID{},
	}, nil)
	syncService := NewSchemaSyncService(mockDB)

	err = syncService.SyncSchema(ctx, dbModel.ID)
	require.NoError(t, err)

	mockDB.AssertExpectations(t)

}

func TestUpdateSchemaCreatesNew(t *testing.T) {
	c := context.Background()
	dataBaseId := uuid.New()
	schemaName := "testSchemaName"

	database := sqlc.Database{
		ID:          dataBaseId,
		DisplayName: "",
		DbName:      "",
		Username:    pgtype.Text{},
		Password:    pgtype.Text{},
		Host:        pgtype.Text{},
		Port:        0,
		Type:        "",
		ProjectID:   uuid.UUID{},
	}

	mockDB := mocks.NewMockQuerier(t)
	mockDB.On("GetDatabaseEntityByDatabaseIdAndTypeAndName", mock.Anything, mock.Anything).Return(sqlc.DatabaseEntity{}, pgx.ErrNoRows)
	mockDB.On("CreateDatabaseEntity", mock.Anything, mock.Anything).Return(sqlc.DatabaseEntity{
		ID:         uuid.New(),
		Name:       schemaName,
		DatabaseID: dataBaseId,
	}, nil)

	syncService := NewSchemaSyncService(mockDB)
	err := syncService.updateSchema(c, schemaName, database)

	require.NoError(t, err)
	mockDB.AssertExpectations(t)

}

func TestUpdateSchemaUpdatesExisting(t *testing.T) {
	c := context.Background()
	dataBaseId := uuid.New()
	schemaName := "testSchemaName"

	database := sqlc.Database{
		ID:          dataBaseId,
		DisplayName: "",
		DbName:      "",
		Username:    pgtype.Text{},
		Password:    pgtype.Text{},
		Host:        pgtype.Text{},
		Port:        0,
		Type:        "",
		ProjectID:   uuid.New(),
	}

	mockDB := mocks.NewMockQuerier(t)
	mockDB.On("GetDatabaseEntityByDatabaseIdAndTypeAndName", mock.Anything, mock.Anything).Return(sqlc.DatabaseEntity{
		ID:         uuid.New(),
		Name:       schemaName,
		DatabaseID: dataBaseId,
	}, nil)

	syncService := NewSchemaSyncService(mockDB)
	err := syncService.updateSchema(c, schemaName, database)

	require.NoError(t, err)
	mockDB.AssertExpectations(t)

}
