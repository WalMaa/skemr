package dbreflect

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/walmaa/skemr-common/models"
)

// Creates a singleton Postgres container for tests
func GetTestPostgres(t *testing.T) *postgres.PostgresContainer {
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
			postgres.WithInitScripts(
				"../../testdata/init_postgres.sql"),
		)

		require.NoError(t, err)
	})
	return pgC
}

// newTestPGConn prepares a temporary Postgres connection for tests.
// It registers a cleanup to close the connection automatically.
func newTestPGConn(t *testing.T) (ctx context.Context, dbConn *PostgresConnector, conn *pgx.Conn) {
	t.Helper()

	ctx = context.Background()

	pgC := GetTestPostgres(t)

	host, err := pgC.Host(ctx)
	require.NoError(t, err)

	port, err := pgC.MappedPort(ctx, "5432")
	require.NoError(t, err)

	// You could parameterize these via env if needed
	dbUser := "user"
	dbPassword := "password"

	hostStr := host
	dbName := "postgres"

	dbModel := models.Database{
		ID:           uuid.New(),
		DisplayName:  "Test Database",
		Username:     &dbUser,
		Password:     &dbPassword,
		Host:         &hostStr,
		DbName:       &dbName,
		Port:         int32(port.Int()),
		DatabaseType: models.Postgres,
		ProjectID:    uuid.New(),
	}

	tmp := NewPostgresConnector(dbModel)
	c, err := tmp.Connect(ctx)
	require.NoError(t, err)
	require.NotNil(t, c)

	// Auto-close after each test
	t.Cleanup(func() {
		_ = c.Close(ctx)
	})

	return ctx, tmp, c
}
