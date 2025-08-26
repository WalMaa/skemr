package dbreflect

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
