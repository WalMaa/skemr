package dbconnector

import (
	"context"
	"errors"

	"github.com/walmaa/skemr/db/sqlc"
)

type DBConnector struct {
	sqlc.Database
}

func newDBConnector(db sqlc.Database) *DBConnector {
	return &DBConnector{
		Database: db,
	}
}

// Connect connects to a database and returns
func (dc *DBConnector) Connect(ctx context.Context) {
	//conn, err := pgx.Connect(ctx)
	//fmt.Print(conn, err)
}

// getConnectionString returns the connection string for the database.
func (dc *DBConnector) getConnectionString() (string, error) {
	host := *dc.Database.Host
	username := *dc.Database.Username
	password := *dc.Database.Password
	port := dc.Database.Port
	if host == "" || port == 0 || dc.Database.DbName == "" {
		return "", errors.New("Missing database connection parameters")

	}

	credentials := ""
	if username == "" || password == "" {
		credentials = username + ":" + password + "@"
	}

	switch dc.Database.Type {
	case "postgres":
		return "postgresql://" + credentials + host + ":" + string(port) + "/" + dc.Database.DbName, nil
	}
	return "", errors.New("DB not supported")
}
