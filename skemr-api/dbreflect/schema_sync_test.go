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

	mockDB.On("GetSchemaByNameAndDatabase", mock.Anything, mock.Anything).Return(sqlc.Schema{
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

	mockDB := mocks.NewMockQuerier(t)
	mockDB.On("GetSchemaByNameAndDatabase", mock.Anything, mock.Anything).Return(sqlc.Schema{}, pgx.ErrNoRows)
	mockDB.On("CreateSchema", mock.Anything, mock.Anything).Return(sqlc.Schema{
		ID:         uuid.New(),
		Name:       schemaName,
		DatabaseID: dataBaseId,
	}, nil)

	syncService := NewSchemaSyncService(mockDB)
	err := syncService.updateSchema(c, schemaName, uuid.New())

	require.NoError(t, err)
	mockDB.AssertExpectations(t)

}

func TestUpdateSchemaUpdatesExisting(t *testing.T) {
	c := context.Background()
	dataBaseId := uuid.New()
	schemaName := "testSchemaName"

	mockDB := mocks.NewMockQuerier(t)
	mockDB.On("GetSchemaByNameAndDatabase", mock.Anything, mock.Anything).Return(sqlc.Schema{
		ID:         uuid.New(),
		Name:       schemaName,
		DatabaseID: dataBaseId,
	}, nil)

	syncService := NewSchemaSyncService(mockDB)
	err := syncService.updateSchema(c, schemaName, uuid.New())

	require.NoError(t, err)
	mockDB.AssertExpectations(t)

}
