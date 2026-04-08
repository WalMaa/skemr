package dbreflect

import (
	"context"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/walmaa/skemr-common/models"
)

var (
	pgOnce sync.Once
	pgC    *postgres.PostgresContainer
)

func TestConnectToPostgres(t *testing.T) {
	pgC := GetTestPostgres(t)
	ctx := context.Background()

	host, err := pgC.Host(ctx)
	require.NoError(t, err)

	port, err := pgC.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dbUser := "user"
	dbPassword := "password"
	dbName := "postgres"

	dbModel := &models.Database{
		ID:           uuid.New(),
		DisplayName:  "Test Database",
		Username:     &dbUser,
		Password:     &dbPassword,
		Host:         &host,
		DatabaseType: models.Postgres,
		DbName:       &dbName,
		Port:         int32(port.Int()),
		ProjectID:    uuid.New(),
	}

	dbConn := NewPostgresConnector(*dbModel)
	conn, err := dbConn.Connect(ctx)
	require.NoError(t, err)
	require.NotNil(t, conn)

	err = conn.Close(ctx)
	require.NoError(t, err)
}

func TestGetConnectionStringWithoutCreds(t *testing.T) {
	host := "localhost"
	dbName := "testDb"
	dbModel := &models.Database{
		ID:           uuid.New(),
		DisplayName:  "Test Database",
		Host:         &host,
		DatabaseType: models.Postgres,
		DbName:       &dbName,
		Port:         5432,
		ProjectID:    uuid.New(),
	}
	dbConn := NewPostgresConnector(*dbModel)
	connStr, err := dbConn.getConnectionString()
	require.NoError(t, err)
	expected := "postgresql://localhost:5432/testDb?sslmode=prefer"
	require.Equal(t, expected, connStr)
}

func TestGetConnectionStringWithoutUsername(t *testing.T) {
	host := "localhost"
	username := "user"
	password := "password"
	dbName := "postgres"
	dbModel := &models.Database{
		ID:           uuid.New(),
		DisplayName:  "Test Database",
		Host:         &host,
		Username:     &username,
		Password:     &password,
		DatabaseType: models.Postgres,
		DbName:       &dbName,
		Port:         5432,
		ProjectID:    uuid.New(),
	}
	dbConn := NewPostgresConnector(*dbModel)
	connStr, err := dbConn.getConnectionString()
	require.NoError(t, err)
	expected := "postgresql://user:password@localhost:5432/postgres?sslmode=prefer"
	require.Equal(t, expected, connStr)
}

func TestGetConnectionStringWithoutPass(t *testing.T) {
	host := "localhost"
	username := "user"
	dbName := "postgres"
	dbModel := &models.Database{
		ID:           uuid.New(),
		DisplayName:  "Test Database",
		Host:         &host,
		Username:     &username,
		DatabaseType: models.Postgres,
		DbName:       &dbName,
		Port:         5432,
		ProjectID:    uuid.New(),
	}
	dbConn := NewPostgresConnector(*dbModel)
	connStr, err := dbConn.getConnectionString()
	require.NoError(t, err)
	expected := "postgresql://localhost:5432/postgres?sslmode=prefer"
	require.Equal(t, expected, connStr)
}

func TestGetConnectionStringWithCreds(t *testing.T) {
	host := "localhost"
	username := "user"
	password := "password"
	dbName := "postgres"
	dbModel := &models.Database{
		ID:           uuid.New(),
		DisplayName:  "Test Database",
		Host:         &host,
		Username:     &username,
		Password:     &password,
		DatabaseType: models.Postgres,
		DbName:       &dbName,
		Port:         5432,
		ProjectID:    uuid.New(),
	}
	dbConn := NewPostgresConnector(*dbModel)
	connStr, err := dbConn.getConnectionString()
	require.NoError(t, err)
	expected := "postgresql://user:password@localhost:5432/postgres?sslmode=prefer"
	require.Equal(t, expected, connStr)
}

func TestGetConnectionWithSslMode(t *testing.T) {
	host := "localhost"
	username := "user"
	password := "password"
	dbName := "postgres"
	sslMode := "require"
	dbModel := &models.Database{
		ID:           uuid.New(),
		DisplayName:  "Test Database",
		Host:         &host,
		Username:     &username,
		Password:     &password,
		DatabaseType: models.Postgres,
		DbName:       &dbName,
		Port:         5432,
		SslMode:      sslMode,
		ProjectID:    uuid.New(),
	}
	dbConn := NewPostgresConnector(*dbModel)
	connStr, err := dbConn.getConnectionString()
	require.NoError(t, err)
	expected := "postgresql://user:password@localhost:5432/postgres?sslmode=require"
	require.Equal(t, expected, connStr)
}

func TestGetConnectionWithVerifyFullSslMode(t *testing.T) {
	host := "localhost"
	username := "user"
	password := "password"
	dbName := "postgres"
	sslMode := "verify-full"
	dbModel := &models.Database{
		ID:           uuid.New(),
		DisplayName:  "Test Database",
		Host:         &host,
		Username:     &username,
		Password:     &password,
		DatabaseType: models.Postgres,
		DbName:       &dbName,
		Port:         5432,
		SslMode:      sslMode,
		ProjectID:    uuid.New(),
	}
	dbConn := NewPostgresConnector(*dbModel)
	connStr, err := dbConn.getConnectionString()
	require.NoError(t, err)
	expected := "postgresql://user:password@localhost:5432/postgres?sslmode=verify-full"
	require.Equal(t, expected, connStr)
}
