package dbreflect

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/walmaa/skemr/db/sqlc"
	"github.com/walmaa/skemr/errormsg"
)

type SchemaSyncService struct {
	db sqlc.Querier
}

func NewSchemaSyncService(db sqlc.Querier) *SchemaSyncService {
	return &SchemaSyncService{db: db}
}

func (s *SchemaSyncService) SyncSchema(c context.Context, databaseID uuid.UUID) error {
	database, err := s.db.GetDatabase(c, databaseID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errormsg.ErrDatabaseNotFound
		}
		return fmt.Errorf("error getting database: %w", err)
	}

	postgresConn := NewPostgresConnector(database)

	conn, err := postgresConn.Connect(c)

	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			fmt.Printf("error closing connection: %v\n", err)
		}
	}(conn, c)

	// Get all schemas in the database
	schemas, err := postgresConn.GetSchemas(c, conn)
	if err != nil {
		return fmt.Errorf("error getting schemas: %w", err)
	}

	// For each schema, get tables and columns
	for _, schema := range schemas {
		err := s.updateSchema(c, schema, databaseID)
		if err != nil {
			return err
		}
		tables, err := postgresConn.GetTablesInSchema(c, conn, schema)
		if err != nil {
			return fmt.Errorf("error getting tables in schema %s: %w", schema, err)
		}
		for _, tableRef := range tables {
			columns, err := postgresConn.ListColumnsInTable(c, conn, tableRef)
			if err != nil {
				return fmt.Errorf("error getting columns in table %s.%s: %w", schema, tableRef.Table, err)
			}
			fmt.Printf("Schema: %s, Table: %s, Columns: %+v\n", schema, tableRef.Table, columns)
			// Here you would typically update your local representation of the schema
			// For example, you might want to upsert the schema, tables, and columns into your database
			// This part is omitted for brevity
		}
	}

	return nil
}

func (s *SchemaSyncService) updateSchema(c context.Context, schemaName string, databaseId uuid.UUID) error {
	args := sqlc.GetSchemaByNameAndDatabaseParams{Name: schemaName, DatabaseID: databaseId}
	schema, err := s.db.GetSchemaByNameAndDatabase(c, args)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("error getting schema", "error", err.Error())
		return err
	}

	// If that schema does not exist yet, save it
	if errors.Is(err, pgx.ErrNoRows) {
		args := sqlc.CreateSchemaParams{
			DatabaseID: databaseId,
			Name:       schemaName,
		}
		schema, err := s.db.CreateSchema(c, args)
		if err != nil {
			slog.Error("error creating schema", "error", err)
			return err
		}
		slog.Info("schema created", "schema", schema)
		return nil
	}

	// If schema does exist, update if needed. Currently no-op
	slog.Info("schema exists", "schema", schema)

	return nil

}

func (s *SchemaSyncService) UpdateTable(c context.Context, databaseId uuid.UUID, tableRef TableRef) error {
	return nil
}
