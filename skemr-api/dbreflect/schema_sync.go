package dbreflect

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/errormsg"
	"github.com/walmaa/skemr-common/models"
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
		schema, err := s.updateSchema(c, schema, database)
		if err != nil {
			return err
		}
		tables, err := postgresConn.GetTablesInSchema(c, conn, schema.Name)
		if err != nil {
			return fmt.Errorf("error getting tables in schema %s: %w", schema, err)
		}
		for _, tableRef := range tables {
			table, err := s.UpdateTable()
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

func (s *SchemaSyncService) updateSchema(c context.Context, schemaName string, database sqlc.Database) (sqlc.DatabaseEntity, error) {
	args := sqlc.GetDatabaseEntityByDatabaseIdAndTypeAndNameParams{
		DatabaseID: database.ID,
		Type:       sqlc.DatabaseEntityType(models.DatabaseEntityTypeSchema),
		Name:       schemaName,
	}
	schema, err := s.db.GetDatabaseEntityByDatabaseIdAndTypeAndName(c, args)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("error getting schema", "error", err.Error())
		return sqlc.DatabaseEntity{}, err
	}

	// If that schema does not exist yet, save it
	if errors.Is(err, pgx.ErrNoRows) {
		args := sqlc.CreateDatabaseEntityParams{
			ProjectID:  database.ProjectID,
			DatabaseID: database.ID,
			Name:       schemaName,
		}
		schema, err := s.db.CreateDatabaseEntity(c, args)
		if err != nil {
			slog.Error("error creating schema", "error", err)
			return sqlc.DatabaseEntity{}, err
		}
		slog.Info("schema created", "schema", schema.Name)
		return schema, err
	}

	// If schema does exist, update if needed. Currently no-op
	slog.Info("schema exists", "schema", schema)

	return schema, err

}

func (s *SchemaSyncService) UpdateTable(c context.Context, tableRef TableRef, database sqlc.Database) (sqlc.DatabaseEntity, error) {
	return sqlc.DatabaseEntity{}, nil
}
