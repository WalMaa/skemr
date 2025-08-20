package dbconnector

import (
	"context"
	"errors"
	"log/slog"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/walmaa/skemr/db/sqlc"
)

type DBConnector struct {
	sqlc.Database
}

func NewDBConnector(db sqlc.Database) *DBConnector {
	return &DBConnector{
		Database: db,
	}
}

// Connect connects to a database and returns a pgx.Conn instance.
func (dc *DBConnector) Connect(ctx context.Context) (*pgx.Conn, error) {
	connStr, err := dc.getConnectionString()
	if err != nil {
		slog.Error("Error getting connection string", "err", err)
		return nil, err
	}

	conn, err := pgx.Connect(ctx, connStr)

	if err != nil {
		slog.Error("Error connecting to database", "connection_string", connStr, "err", err)
		return nil, err
	}
	slog.Info("Connected to database", "connection_string", connStr)
	return conn, nil
}

// getConnectionString returns the connection string for the database.
func (dc *DBConnector) getConnectionString() (string, error) {
	host := *dc.Database.Host
	port := dc.Database.Port
	if host == "" || port == 0 || dc.Database.DbName == "" {
		return "", errors.New("Missing database connection parameters")

	}

	credentials := ""
	if dc.Database.Username != nil && dc.Database.Password != nil {
		username := *dc.Database.Username
		password := *dc.Database.Password
		credentials = username + ":" + password + "@"
	}

	switch dc.Database.Type {
	case "postgres":
		return "postgresql://" + credentials + host + ":" + strconv.Itoa(int(port)) + "/" + dc.Database.DbName, nil
	}
	return "", errors.New("DB not supported")
}
