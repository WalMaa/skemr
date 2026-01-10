package dbreflect

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/mapper"
	"github.com/walmaa/skemr-api/tasks"
	"github.com/walmaa/skemr-common/models"
)

type SchemaSyncService struct {
	db sqlc.Querier
}

func NewSchemaSyncService(db sqlc.Querier) *SchemaSyncService {
	return &SchemaSyncService{db: db}
}

func (s *SchemaSyncService) ProcessSyncTask(c context.Context, t *asynq.Task) error {
	slog.Info("Starting database sync task")
	var p tasks.DatabaseSyncPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	slog.Info("Syncing database ", slog.String("databaseId", p.DatabaseID.String()))

	database, err := s.db.GetDatabase(c, p.DatabaseID)

	if err != nil {
		slog.Error("Could not get database for database sync", slog.String("databaseId", p.DatabaseID.String()))
		return err
	}
	err = s.SyncSchema(c, mapper.ToDomainDatabase(database))
	if err != nil {
		return err
	}

	return nil
}

func (s *SchemaSyncService) SyncSchema(c context.Context, database models.Database) error {

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
			return fmt.Errorf("error getting tables in schema %q: %w", schema.Name, err)
		}
		for _, tableRef := range tables {
			table, err := s.UpdateTable(c, tableRef, database, schema.ID)
			columns, err := postgresConn.ListColumnsInTable(c, conn, tableRef)
			if err != nil {
				return fmt.Errorf("error getting columns in table %q.%q: %w", schema.Name, tableRef.Name, err)
			}
			for _, column := range columns {
				_, err := s.UpdateColumn(c, column, database, table.ID)
				if err != nil {
					return fmt.Errorf("Error updating tables: %w", err)
				}
			}

		}
	}

	return nil
}

func (s *SchemaSyncService) updateSchema(c context.Context, schemaName string, database models.Database) (sqlc.DatabaseEntity, error) {
	args := sqlc.GetDatabaseEntityByDatabaseIdAndTypeAndNameParams{
		DatabaseID: database.ID,
		EntityType: sqlc.DatabaseEntityTypeSchema,
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
			EntityType: sqlc.DatabaseEntityTypeSchema,
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

func (s *SchemaSyncService) UpdateTable(c context.Context, tableRef TableRef, database models.Database, schemaId uuid.UUID) (sqlc.DatabaseEntity, error) {
	args := sqlc.GetDatabaseEntityByDatabaseIdAndTypeAndNameParams{
		DatabaseID: database.ID,
		EntityType: sqlc.DatabaseEntityTypeTable,
		Name:       tableRef.Name,
	}
	table, err := s.db.GetDatabaseEntityByDatabaseIdAndTypeAndName(c, args)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("error getting schema", "error", err.Error())
		return sqlc.DatabaseEntity{}, err
	}

	// If that schema does not exist yet, save it
	if errors.Is(err, pgx.ErrNoRows) {
		args := sqlc.CreateDatabaseEntityParams{
			ProjectID:  database.ProjectID,
			EntityType: sqlc.DatabaseEntityTypeTable,
			ParentID:   &schemaId,
			DatabaseID: database.ID,
			Name:       tableRef.Name,
		}
		table, err := s.db.CreateDatabaseEntity(c, args)
		if err != nil {
			slog.Error("error creating schema", "error", err)
			return sqlc.DatabaseEntity{}, err
		}
		slog.Info("schema created", "schema", table.Name)
		return table, err
	}

	// If table does exist, update if needed. Currently no-op
	slog.Info("table exists", "schema", table)

	return table, err
}

func (s *SchemaSyncService) UpdateColumn(c context.Context, columnRef ColumnRef, database models.Database, tableId uuid.UUID) (sqlc.DatabaseEntity, error) {
	args := sqlc.GetDatabaseEntityByDatabaseIdAndTypeAndNameParams{
		DatabaseID: database.ID,
		EntityType: sqlc.DatabaseEntityTypeColumn,
		Name:       columnRef.Name,
	}
	column, err := s.db.GetDatabaseEntityByDatabaseIdAndTypeAndName(c, args)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("error getting schema", "error", err.Error())
		return sqlc.DatabaseEntity{}, err
	}

	// If that schema does not exist yet, save it
	if errors.Is(err, pgx.ErrNoRows) {
		args := sqlc.CreateDatabaseEntityParams{
			ProjectID:  database.ProjectID,
			EntityType: sqlc.DatabaseEntityTypeColumn,
			ParentID:   &tableId,
			DatabaseID: database.ID,
			Name:       columnRef.Name,
		}
		column, err := s.db.CreateDatabaseEntity(c, args)
		if err != nil {
			slog.Error("error creating schema", "error", err)
			return sqlc.DatabaseEntity{}, err
		}
		slog.Info("schema created", "schema", column.Name)
		return column, err
	}

	// If column does exist, update if needed. Currently no-op
	slog.Info("column exists", "schema", column)

	return column, err
}
