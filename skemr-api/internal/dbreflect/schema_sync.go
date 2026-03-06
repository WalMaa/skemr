package dbreflect

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/mapper"
	"github.com/walmaa/skemr-api/internal/tasks"
	"github.com/walmaa/skemr-common/models"
)

type SchemaSyncService struct {
	db               sqlc.Querier
	connectorFactory func(database models.Database) DatabaseConnector
}

type ColumnAttributes struct {
	DataType  string  `json:"dataType"`
	Default   *string `json:"default"`
	Nullable  string  `json:"nullable"`  // YES or NO
	Updatable string  `json:"updatable"` // YES or NO
}

func NewSchemaSyncService(db sqlc.Querier, connectorFactory func(database models.Database) DatabaseConnector) *SchemaSyncService {
	return &SchemaSyncService{db: db, connectorFactory: connectorFactory}
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
		slog.Error("Error syncing database schema", "databaseId", p.DatabaseID, "error", err)
		_, err := s.db.UpdateDatabaseSyncFail(c, sqlc.UpdateDatabaseSyncFailParams{
			SyncError:  pgtype.Text{String: err.Error(), Valid: true},
			SyncedAt:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
			DatabaseID: p.DatabaseID,
		})
		if err != nil {
			slog.Error("error updating database sync fail", "error", err)
			// Do not return this error as we want to keep the original sync error for debugging.
		}

		return err
	}

	// Update database object with last synced time
	_, err = s.db.UpdateDatabaseSyncedAt(c, sqlc.UpdateDatabaseSyncedAtParams{
		SyncedAt:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
		DatabaseID: database.ID,
	},
	)
	if err != nil {
		slog.Error("error updating database synced at", "error", err)
		return err
	}

	return nil
}

func (s *SchemaSyncService) SyncSchema(c context.Context, database models.Database) error {

	connector := s.connectorFactory(database)
	conn, err := connector.Connect(c)

	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		if conn == nil {
			slog.Warn("No connection to close")
			return
		}
		err := conn.Close(ctx)
		if err != nil {
			fmt.Printf("error closing connection: %v\n", err)
		}
	}(conn, c)

	// Get all saved entities for the database. Any entities that are not found in the new schema will be marked as deleted at the end.
	savedEntities, err := s.db.GetDatabaseEntitiesByDatabaseId(c, database.ID)
	if err != nil {
		return fmt.Errorf("error getting saved database entities: %w", err)
	}
	// Create a map of current entity ids to easily check which entities are still present in the new schema after the sync.
	currentEntityIds := make([]uuid.UUID, 0, len(savedEntities))

	// Get all schemas in the database
	schemas, err := connector.GetSchemas(c, conn)
	if err != nil {
		return fmt.Errorf("error getting schemas: %w", err)
	}

	// For each schema, get tables and columns
	for _, schema := range schemas {
		schema, err := s.updateSchema(c, schema, database)
		if err != nil {
			return err
		}
		// Add schema ID to current entity ids
		slog.Debug("Adding schema to current entity ids", "schemaName", schema.Name, "schemaId", schema.ID)
		currentEntityIds = append(currentEntityIds, schema.ID)
		tables, err := connector.GetTablesInSchema(c, conn, schema.Name)
		if err != nil {
			return fmt.Errorf("error getting tables in schema %q: %w", schema.Name, err)
		}
		for _, tableRef := range tables {
			table, err := s.UpdateTable(c, tableRef, database, schema.ID)

			if err != nil {
				return fmt.Errorf("error updating tables: %w", err)
			}
			// Add table ID to current entity ids
			slog.Debug("Adding table to current entity ids", "tableName", table.Name, "tableId", table.ID)
			currentEntityIds = append(currentEntityIds, table.ID)
			columns, err := connector.ListColumnsInTable(c, conn, tableRef)
			if err != nil {
				return fmt.Errorf("error getting columns in table %q.%q: %w", schema.Name, tableRef.Name, err)
			}
			for _, column := range columns {
				column, err := s.SyncColumn(c, column, database, table.ID)
				if err != nil {
					return fmt.Errorf("Error updating column: %w", err)
				}

				// Add column ID to current entity ids
				slog.Debug("Adding column to current entity ids", "columnName", column.Name, "columnId", column.ID)
				currentEntityIds = append(currentEntityIds, column.ID)
			}

		}
	}

	// Mark any entities that were not found in the new schema as deleted
	for _, entity := range savedEntities {
		if !slices.Contains(currentEntityIds, entity.ID) {
			err := s.markEntityAsDeleted(c, entity.ID)
			if err != nil {
				slog.Error("Error marking entity as deleted", "entityId", entity.ID, "error", err)
				// Do not return error as we want to continue marking other entities as deleted
			}
		}
	}

	return nil
}

// updateSchema checks if a schema with the given name exists for the database.
// If it does not exist, it creates a new schema entity.
// If it does exist, it currently does nothing but can be extended to update schema attributes if needed.
func (s *SchemaSyncService) updateSchema(c context.Context, schemaName string, database models.Database) (sqlc.DatabaseEntity, error) {
	args := sqlc.GetDatabaseEntityByDatabaseIdAndTypeAndParentAndNameParams{
		DatabaseID: database.ID,
		EntityType: sqlc.DatabaseEntityTypeSchema,
		Name:       schemaName,
	}
	schema, err := s.db.GetDatabaseEntityByDatabaseIdAndTypeAndParentAndName(c, args)
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

	// TODO: If schema does exist, update if needed. Currently no-op
	slog.Info("schema exists", "schema", schema)

	return schema, err

}

// markEntityAsDeleted sets the status of the entity to "deleted".
// This is used for entities that were not found in the new schema during sync.
// We want to keep these entities in the database to show how the schema has changed.
func (s *SchemaSyncService) markEntityAsDeleted(c context.Context, entityId uuid.UUID) error {
	slog.Debug("Marking entity as deleted", "entityId", entityId)
	err := s.db.UpdateDatabaseEntityAsDeleted(c, entityId)
	if err != nil {
		slog.Error("error marking entity as deleted", "error", err)
		return err
	}
	return nil
}

func (s *SchemaSyncService) UpdateTable(c context.Context, tableRef TableRef, database models.Database, schemaId uuid.UUID) (sqlc.DatabaseEntity, error) {
	args := sqlc.GetDatabaseEntityByDatabaseIdAndTypeAndParentAndNameParams{
		DatabaseID: database.ID,
		ParentID:   &schemaId,
		EntityType: sqlc.DatabaseEntityTypeTable,
		Name:       tableRef.Name,
	}
	table, err := s.db.GetDatabaseEntityByDatabaseIdAndTypeAndParentAndName(c, args)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("error getting schema", "error", err.Error())
		return sqlc.DatabaseEntity{}, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		// If that schema does not exist by name, check by fingerprint to see if it is the same table with an updated name.

		fingerprint := GenerateTableFingerprint(tableRef, schemaId)

		table, err = s.db.GetDatabaseEntityByFingerprint(c, sqlc.GetDatabaseEntityByFingerprintParams{
			DatabaseID: database.ID,
			Fingerprint: pgtype.Text{
				String: fingerprint,
				Valid:  true,
			},
		})

		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			slog.Error("error getting table by fingerprint", "error", err.Error())
			return sqlc.DatabaseEntity{}, err
		}

		// If found by fingerprint, update the name to the new name.
		if !errors.Is(err, pgx.ErrNoRows) {
			slog.Debug("Table found by fingerprint, updating name", "oldName", table.Name, "newName", tableRef.Name)

			table, err = s.db.UpdateDatabaseEntityName(c, sqlc.UpdateDatabaseEntityNameParams{
				ID:   table.ID,
				Name: tableRef.Name,
			})
			if err != nil {
				slog.Error("error updating table name", "error", err.Error())
				return sqlc.DatabaseEntity{}, err
			}
			slog.Info("Table renamed", "oldName", table.Name, "newName", tableRef.Name)
			table.Name = tableRef.Name
			return table, nil
		}

		args := sqlc.CreateDatabaseEntityParams{
			ProjectID:  database.ProjectID,
			EntityType: sqlc.DatabaseEntityTypeTable,
			ParentID:   &schemaId,
			DatabaseID: database.ID,
			Name:       tableRef.Name,
			Fingerprint: pgtype.Text{
				String: GenerateTableFingerprint(tableRef, schemaId),
				Valid:  true,
			},
		}
		table, err := s.db.CreateDatabaseEntity(c, args)
		if err != nil {
			slog.Error("error creating schema", "error", err)
			return sqlc.DatabaseEntity{}, err
		}
		slog.Info("Table created", "schema", table.Name)
		return table, err
	}

	// TODO: If table does exist, update if needed. Currently no-op
	slog.Info("table exists", "table", table.Name)

	return table, err
}

// SyncColumn syncs a database column to a database entity.
func (s *SchemaSyncService) SyncColumn(c context.Context, columnRef ColumnRef, database models.Database, tableId uuid.UUID) (sqlc.DatabaseEntity, error) {
	args := sqlc.GetDatabaseEntityByDatabaseIdAndTypeAndParentAndNameParams{
		DatabaseID: database.ID,
		ParentID:   &tableId,
		EntityType: sqlc.DatabaseEntityTypeColumn,
		Name:       columnRef.Name,
	}
	column, err := s.db.GetDatabaseEntityByDatabaseIdAndTypeAndParentAndName(c, args)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("error getting schema", "error", err.Error())
		return sqlc.DatabaseEntity{}, err
	}

	columnAttributes := ColumnAttributes{
		DataType:  columnRef.DataType,
		Default:   columnRef.Default,
		Nullable:  columnRef.Nullable,
		Updatable: columnRef.Updatable,
	}
	attributesJson, jsonError := json.Marshal(columnAttributes)
	if jsonError != nil {
		slog.Error("error marshalling column attributes", "error", err)
		return sqlc.DatabaseEntity{}, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		// If not found by name, check by fingerprint to see if it is the same column with an updated name.
		fingerprint := GenerateColumnFingerprint(columnRef, tableId)

		column, err = s.db.GetDatabaseEntityByFingerprint(c, sqlc.GetDatabaseEntityByFingerprintParams{
			DatabaseID: database.ID,
			Fingerprint: pgtype.Text{
				String: fingerprint,
				Valid:  true,
			},
		})

		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			slog.Error("error getting column by fingerprint", "error", err.Error())
			return sqlc.DatabaseEntity{}, err
		}

		// If found by fingerprint, update the name to the new name.
		if !errors.Is(err, pgx.ErrNoRows) {
			slog.Debug("Column found by fingerprint, updating name", "oldName", column.Name, "newName", columnRef.Name)

			column, err = s.db.UpdateDatabaseEntityName(c, sqlc.UpdateDatabaseEntityNameParams{
				ID:   column.ID,
				Name: columnRef.Name,
			})
			if err != nil {
				slog.Error("error updating column name", "error", err.Error())
				return sqlc.DatabaseEntity{}, err
			}
			slog.Info("Column renamed", "oldName", column.Name, "newName", columnRef.Name)
			column.Name = columnRef.Name
			return column, nil
		}

		// If that column does not exist yet, save it
		args := sqlc.CreateDatabaseEntityParams{
			ProjectID:  database.ProjectID,
			EntityType: sqlc.DatabaseEntityTypeColumn,
			ParentID:   &tableId,
			DatabaseID: database.ID,
			Name:       columnRef.Name,
			Attributes: attributesJson,
			Fingerprint: pgtype.Text{
				String: GenerateColumnFingerprint(columnRef, tableId),
				Valid:  true,
			},
		}
		column, err := s.db.CreateDatabaseEntity(c, args)
		if err != nil {
			slog.Error("error creating schema", "error", err)
			return sqlc.DatabaseEntity{}, err
		}
		slog.Info("Column created", "name", column.Name)
		return column, err
	}

	// If column does exist, update if needed. Currently no-op
	slog.Info("column exists", "column", column.Name)

	return column, err
}
