package dbconnector

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/walmaa/skemr/db/sqlc"
)

type DBConnector struct {
	sqlc.Database
}

type TableRef struct {
	Schema string
	Table  string
}

type ColumnRef struct {
	Table      TableRef
	ColumnName string
	DataType   string
}

func NewDBConnector(db sqlc.Database) *DBConnector {
	return &DBConnector{
		Database: db,
	}
}

func (dc *DBConnector) ListColumnsInTable(ctx context.Context, conn *pgx.Conn, tableRef TableRef) ([]ColumnRef, error) {
	rows, err := conn.Query(ctx, "SELECT column_name, data_type FROM information_schema.columns WHERE table_schema=$1 AND table_name=$2", tableRef.Schema, tableRef.Table)
	if err != nil {
		slog.Error("Error querying columns", "schema", tableRef.Schema, "table", tableRef.Table, "err", err)
		return nil, err
	}
	defer rows.Close()
	var columns []ColumnRef
	for rows.Next() {
		columnRef := ColumnRef{Table: tableRef}
		if err := rows.Scan(&columnRef.ColumnName, &columnRef.DataType); err != nil {
			slog.Error("Error scanning column name", "err", err)
			return nil, err
		}
		columns = append(columns, columnRef)
	}
	return columns, rows.Err()
}

func (dc *DBConnector) ListTablesInSchema(ctx context.Context, conn *pgx.Conn, schema string) ([]TableRef, error) {
	rows, err := conn.Query(ctx, "SELECT table_schema, table_name FROM information_schema.tables WHERE table_schema=$1", schema)
	if err != nil {
		slog.Error("Error querying tables", "schema", schema, "err", err)
		return nil, err
	}
	defer rows.Close()
	var tables []TableRef
	for rows.Next() {
		var tableRef TableRef
		if err := rows.Scan(&tableRef.Schema, &tableRef.Table); err != nil {
			slog.Error("Error scanning table name", "err", err)
			return nil, err
		}
		tables = append(tables, tableRef)
	}
	return tables, rows.Err()
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

func (dc *DBConnector) Disconnect(ctx context.Context, conn *pgx.Conn) error {
	if conn == nil {
		return errors.New("No connection to close")
	}
	return conn.Close(ctx)
}

func (dc *DBConnector) TestConnection(ctx context.Context) error {
	conn, err := dc.Connect(ctx)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer func() {
		if cerr := dc.Disconnect(ctx, conn); cerr != nil {
			slog.Error("disconnect failed", "err", cerr)
		}
	}()

	if pingErr := conn.Ping(ctx); pingErr != nil {
		slog.Error("ping failed", "err", pingErr)
		return fmt.Errorf("ping: %w", pingErr)
	}

	slog.Info("Successfully connected and pinged the database")
	return nil
}

// getConnectionString returns the connection string for the database.
func (dc *DBConnector) getConnectionString() (string, error) {
	host := dc.Database.Host
	port := dc.Database.Port
	if !host.Valid || port == 0 || dc.Database.DbName == "" {
		return "", errors.New("Missing database connection parameters")
	}

	credentials := ""
	if dc.Database.Username.Valid && dc.Database.Password.Valid {
		username := dc.Database.Username.String
		password := dc.Database.Password.String
		credentials = username + ":" + password + "@"
	}

	switch dc.Database.Type {
	case "postgres":
		return "postgresql://" + credentials + host.String + ":" + strconv.Itoa(int(port)) + "/" + dc.Database.DbName, nil
	}
	return "", errors.New("DB not supported")
}
