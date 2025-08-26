package dbreflect

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/walmaa/skemr/db/sqlc"
	"github.com/walmaa/skemr/test/mocks"
)

func TestSchemaSync(t *testing.T) {
	ctx := context.Background()
	pgC := GetTestPostgres(t)

	host, err := pgC.Host(ctx)
	require.NoError(t, err)

	port, err := pgC.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dbModel := &sqlc.Database{
		ID:          uuid.New(),
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
	synService := NewSchemaSyncService(mockDB)

	err = synService.SyncSchema(ctx, dbModel.ID)
	require.NoError(t, err)

	mockDB.AssertExpectations(t)

}
