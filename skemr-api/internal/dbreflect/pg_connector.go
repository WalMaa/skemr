package dbreflect

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/walmaa/skemr-common/models"
)

type PostgresConnector struct {
	models.Database
}

type TableRef struct {
	Schema string // The parent Schema
	Name   string // The name of the table itself
}

type ColumnRef struct {
	Name      string
	DataType  string
	Default   *string
	Nullable  bool
	Updatable bool
}

func NewPostgresConnector(db models.Database) *PostgresConnector {
	return &PostgresConnector{
		Database: db,
	}
}

func (dc *PostgresConnector) ListColumnsInTable(ctx context.Context, conn *pgx.Conn, tableRef TableRef) ([]ColumnRef, error) {
	rows, err := conn.Query(ctx, "SELECT column_name, data_type, column_default, is_nullable, is_updatable FROM information_schema.columns WHERE table_schema=$1 AND table_name=$2", tableRef.Schema, tableRef.Name)
	if err != nil {
		slog.Error("Error querying columns", "schema", tableRef.Schema, "table", tableRef.Name, "err", err)
		return nil, err
	}
	defer rows.Close()
	var columns []ColumnRef
	for rows.Next() {
		var columnRef ColumnRef
		if err := rows.Scan(&columnRef.Name, &columnRef.DataType, &columnRef.Default, &columnRef.Nullable, &columnRef.Updatable); err != nil {
			slog.Error("Error scanning column name", "err", err)
			return nil, err
		}
		columns = append(columns, columnRef)
	}
	return columns, rows.Err()
}

func (dc *PostgresConnector) GetTablesInSchema(ctx context.Context, conn *pgx.Conn, schema string) ([]TableRef, error) {
	rows, err := conn.Query(ctx, "SELECT table_schema, table_name FROM information_schema.tables WHERE table_schema=$1", schema)
	if err != nil {
		slog.Error("Error querying tables", "schema", schema, "err", err)
		return nil, err
	}
	defer rows.Close()
	var tables []TableRef
	for rows.Next() {
		var tableRef TableRef
		if err := rows.Scan(&tableRef.Schema, &tableRef.Name); err != nil {
			slog.Error("Error scanning table name", "err", err)
			return nil, err
		}
		tables = append(tables, tableRef)
	}
	return tables, rows.Err()
}

func (dc *PostgresConnector) GetSchemas(ctx context.Context, conn *pgx.Conn) ([]string, error) {
	rows, err := conn.Query(ctx, "SELECT schema_name"+
		" FROM information_schema.schemata"+
		" WHERE schema_name NOT IN ('pg_catalog', 'information_schema')"+
		" AND schema_name NOT LIKE 'pg_toast%'"+
		" AND schema_name NOT LIKE 'pg_temp%';")
	if err != nil {
		slog.Error("Error querying schemas", "err", err)
		return nil, err
	}

	defer rows.Close()
	var schemas []string
	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			slog.Error("Error scanning schema name", "err", err)
			return nil, err
		}
		schemas = append(schemas, schema)
	}
	return schemas, rows.Err()
}

// Connect connects to a database and returns a pgx.Conn instance.
func (dc *PostgresConnector) Connect(ctx context.Context) (*pgx.Conn, error) {
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

func (dc *PostgresConnector) Disconnect(ctx context.Context, conn *pgx.Conn) error {
	if conn == nil {
		return errors.New("No connection to close")
	}
	return conn.Close(ctx)
}

func (dc *PostgresConnector) TestConnection(ctx context.Context) error {
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
func (dc *PostgresConnector) getConnectionString() (string, error) {
	host := dc.Database.Host
	port := dc.Database.Port
	if host == nil || port == 0 || dc.DbName == nil {
		return "", errors.New("Missing database connection parameters")
	}

	credentials := ""
	if dc.Database.Username != nil && dc.Database.Password != nil {
		username := *dc.Database.Username
		password := *dc.Database.Password
		credentials = username + ":" + password + "@"
	}

	switch dc.Database.DatabaseType {
	case "postgres":
		return "postgresql://" + credentials + *host + ":" + strconv.Itoa(int(port)) + "/" + *dc.Database.DbName, nil
	}
	return "", errors.New("DB not supported")
}
