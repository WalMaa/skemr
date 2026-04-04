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

const tableDefQuery = `
WITH rels AS ( --- Get all user-defined tables and their OIDs
    SELECT
        t.table_schema,
        t.table_name,
        c.oid AS relid
    FROM information_schema.tables t
             JOIN pg_namespace n
                  ON n.nspname = t.table_schema
             JOIN pg_class c
                  ON c.relnamespace = n.oid
                      AND c.relname = t.table_name
    WHERE t.table_type IN ('BASE TABLE', 'FOREIGN')
      AND t.table_schema = $1
),
     coldefs AS ( -- Get column definitions for each table, including data type, nullability, identity, generated status, and default expressions
         SELECT
             r.table_schema,
             r.table_name,
             r.relid,
             a.attnum,
             lower(format_type(a.atttypid, a.atttypmod)) AS data_type,
             a.attnotnull,
             a.attidentity::text, --- Convert identity column info to text for easier processing
             a.attgenerated::text, --- Convert generated column info to text for easier processing
             pg_get_expr(ad.adbin, ad.adrelid) AS default_expr --- convert expression to text
         FROM rels r
                  JOIN pg_attribute a
                       ON a.attrelid = r.relid
                           AND a.attnum > 0 --- Only user-defined columns (exclude system columns)
                           AND NOT a.attisdropped
                  LEFT JOIN pg_attrdef ad
                            ON ad.adrelid = a.attrelid
                                AND ad.adnum = a.attnum
     ),
     pk_shape AS (
         SELECT
             r.relid,
             COALESCE(string_agg(k.attnum::text, ',' ORDER BY k.ordinality), '') AS pk_shape
         FROM rels r
                  LEFT JOIN pg_constraint pk
                            ON pk.conrelid = r.relid
                                AND pk.contype = 'p'
                  LEFT JOIN LATERAL unnest(pk.conkey) WITH ORDINALITY AS k(attnum, ordinality)
                            ON TRUE
         GROUP BY r.relid
     )
SELECT
    c.table_schema,
    c.table_name,
    string_agg(
            concat_ws(
                    ':',
                    c.data_type, --- include column data type
                    CASE WHEN c.attnotnull THEN 'NO' ELSE 'YES' END, --- include nullability
                    CASE --- determine default type: 'none' if no default, 'sequence' if it's a nextval() default, 'volatile_time' if it's a time-based default, 'generated' if it's a generated column, otherwise just 'default'
                        WHEN c.attgenerated <> '' THEN 'generated'
                        WHEN c.default_expr IS NULL THEN 'none'
                        WHEN c.default_expr ~* '^nextval\(' THEN 'sequence'
                        WHEN c.default_expr ~* '^(now|current_timestamp|transaction_timestamp)\(' THEN 'volatile_time'
                        ELSE 'default'
                        END,
                    COALESCE(NULLIF(c.attidentity, ''), 'none'), --- if an identity column, mark as 'identity', otherwise 'none'
                    COALESCE(NULLIF(c.attgenerated, ''), 'none')
            ),
            ';' ORDER BY c.attnum
    ) AS column_shape,
    pk.pk_shape
FROM coldefs c
         JOIN pk_shape pk
              ON pk.relid = c.relid
GROUP BY c.table_schema, c.table_name, c.relid, pk.pk_shape;`

const schemaRefQuery = `
SELECT string_agg(c.relname, ',' ORDER BY c.relname) AS schema_fingerprint,
       n.nspname AS schema_name
FROM pg_namespace n
         JOIN pg_class c
              ON n.oid = c.relnamespace
WHERE n.nspname NOT IN ('pg_catalog', 'information_schema')
  AND n.nspname NOT LIKE 'pg_%'
  AND c.relkind IN ('r', 'p')
GROUP BY n.nspname;
`

const columnRefQuery = `
SELECT column_name, data_type, column_default, is_nullable, is_updatable, ordinal_position
FROM information_schema.columns
WHERE table_schema = $1 AND table_name = $2
ORDER BY ordinal_position;
`

type TableRef struct {
	Schema      string // The parent Schema
	Name        string // The name of the table itself
	ColumnShape string // A string representation of the column names and data types, used for fingerprinting
	PrimaryKey  string // A string representation of the primary key columns, used for fingerprinting
}

type ColumnRef struct {
	Name            string
	DataType        string
	Default         *string
	Nullable        string // YES or NO
	Updatable       string // YES or NO
	OrdinalPosition int
}

type SchemaRef struct {
	Name        string
	Fingerprint string
}

type DatabaseConnector interface {
	Connect(ctx context.Context) (*pgx.Conn, error)
	Disconnect(ctx context.Context, conn *pgx.Conn) error
	TestConnection(ctx context.Context) error
	GetSchemas(ctx context.Context, conn *pgx.Conn) ([]SchemaRef, error)
	GetTablesInSchema(ctx context.Context, conn *pgx.Conn, schema string) ([]TableRef, error)
	ListColumnsInTable(ctx context.Context, conn *pgx.Conn, tableRef TableRef) ([]ColumnRef, error)
	getConnectionString() (string, error)
}

func NewPostgresConnector(db models.Database) DatabaseConnector {
	return &PostgresConnector{
		Database: db,
	}
}

func (dc *PostgresConnector) ListColumnsInTable(ctx context.Context, conn *pgx.Conn, tableRef TableRef) ([]ColumnRef, error) {
	rows, err := conn.Query(ctx, columnRefQuery, tableRef.Schema, tableRef.Name)
	if err != nil {
		slog.Error("Error querying columns", "schema", tableRef.Schema, "table", tableRef.Name, "err", err)
		return nil, err
	}
	defer rows.Close()
	var columns []ColumnRef
	for rows.Next() {
		var columnRef ColumnRef
		if err := rows.Scan(&columnRef.Name, &columnRef.DataType, &columnRef.Default, &columnRef.Nullable, &columnRef.Updatable, &columnRef.OrdinalPosition); err != nil {
			slog.Error("Error scanning column name", "err", err)
			return nil, err
		}
		columns = append(columns, columnRef)
	}
	return columns, rows.Err()
}

func (dc *PostgresConnector) GetTablesInSchema(ctx context.Context, conn *pgx.Conn, schema string) ([]TableRef, error) {
	rows, err := conn.Query(ctx, tableDefQuery, schema)
	if err != nil {
		slog.Error("Error querying tables", "schema", schema, "err", err)
		return nil, err
	}
	defer rows.Close()
	var tables []TableRef
	for rows.Next() {
		var tableRef TableRef
		if err := rows.Scan(&tableRef.Schema, &tableRef.Name, &tableRef.ColumnShape, &tableRef.PrimaryKey); err != nil {
			slog.Error("Error scanning table name", "err", err)
			return nil, err
		}
		tables = append(tables, tableRef)
	}
	return tables, rows.Err()
}

func (dc *PostgresConnector) GetSchemas(ctx context.Context, conn *pgx.Conn) ([]SchemaRef, error) {
	rows, err := conn.Query(ctx, schemaRefQuery)
	if err != nil {
		slog.Error("Error querying schemas", "err", err)
		return nil, err
	}

	defer rows.Close()
	var schemas []SchemaRef
	for rows.Next() {
		var schemaRef SchemaRef
		if err := rows.Scan(&schemaRef.Fingerprint, &schemaRef.Name); err != nil {
			slog.Error("Error scanning schema name", "err", err)
			return nil, err
		}
		schemas = append(schemas, schemaRef)
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
