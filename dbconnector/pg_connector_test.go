package dbconnector

import (
	"context"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/walmaa/skemr/db/sqlc"
)

var (
	pgOnce sync.Once
	pgC    *postgres.PostgresContainer
)

func getPostgres(t *testing.T) *postgres.PostgresContainer {
	t.Helper()
	var err error
	pgOnce.Do(func() {

		ctx := context.Background()
		dbName := "postgres"
		dbUser := "user"
		dbPassword := "password"

		pgC, err = postgres.Run(ctx,
			"postgres:16-alpine",
			postgres.WithDatabase(dbName),
			postgres.WithUsername(dbUser),
			postgres.WithPassword(dbPassword),
			postgres.BasicWaitStrategies(),
		)

		require.NoError(t, err)
	})
	return pgC
}

func TestConnectToPostgres(t *testing.T) {
	pgC := getPostgres(t)
	ctx := context.Background()

	host, err := pgC.Host(ctx)
	require.NoError(t, err)

	port, err := pgC.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dbUser := "user"
	dbPassword := "password"

	dbModel := &sqlc.Database{
		ID:          uuid.New(),
		DisplayName: "Test Database",
		Username:    &dbUser,
		Password:    &dbPassword,
		Host:        &host,
		Type:        "postgres",
		DbName:      "postgres",
		Port:        int32(port.Int()),
		ProjectID:   uuid.New(),
	}

	dbConn := NewDBConnector(*dbModel)
	conn, err := dbConn.Connect(ctx)
	require.NoError(t, err)
	require.NotNil(t, conn)

	err = conn.Close(ctx)
	require.NoError(t, err)
}

func TestGetConnectionStringWithoutCreds(t *testing.T) {
	host := "localhost"
	dbModel := &sqlc.Database{
		ID:          uuid.New(),
		DisplayName: "Test Database",
		Host:        &host,
		Type:        "postgres",
		DbName:      "testdb",
		Port:        5432,
		ProjectID:   uuid.New(),
	}
	dbConn := NewDBConnector(*dbModel)
	connStr, err := dbConn.getConnectionString()
	require.NoError(t, err)
	expected := "postgresql://localhost:5432/testdb"
	require.Equal(t, expected, connStr)
}

func TestGetConnectionStringWithoutUsername(t *testing.T) {
	host := "localhost"
	username := "user"
	password := "password"
	dbModel := &sqlc.Database{
		ID:          uuid.New(),
		DisplayName: "Test Database",
		Host:        &host,
		Username:    &username,
		Password:    &password,
		Type:        "postgres",
		DbName:      "testdb",
		Port:        5432,
		ProjectID:   uuid.New(),
	}
	dbConn := NewDBConnector(*dbModel)
	connStr, err := dbConn.getConnectionString()
	require.NoError(t, err)
	expected := "postgresql://user:password@localhost:5432/testdb"
	require.Equal(t, expected, connStr)
}

func TestGetConnectionStringWithoutPass(t *testing.T) {
	host := "localhost"
	username := "user"
	dbModel := &sqlc.Database{
		ID:          uuid.New(),
		DisplayName: "Test Database",
		Host:        &host,
		Username:    &username,
		Type:        "postgres",
		DbName:      "testdb",
		Port:        5432,
		ProjectID:   uuid.New(),
	}
	dbConn := NewDBConnector(*dbModel)
	connStr, err := dbConn.getConnectionString()
	require.NoError(t, err)
	expected := "postgresql://localhost:5432/testdb"
	require.Equal(t, expected, connStr)
}

func TestGetConnectionStringWithCreds(t *testing.T) {
	host := "localhost"
	username := "user"
	password := "password"
	dbModel := &sqlc.Database{
		ID:          uuid.New(),
		DisplayName: "Test Database",
		Host:        &host,
		Username:    &username,
		Password:    &password,
		Type:        "postgres",
		DbName:      "testdb",
		Port:        5432,
		ProjectID:   uuid.New(),
	}
	dbConn := NewDBConnector(*dbModel)
	connStr, err := dbConn.getConnectionString()
	require.NoError(t, err)
	expected := "postgresql://user:password@localhost:5432/testdb"
	require.Equal(t, expected, connStr)
}
