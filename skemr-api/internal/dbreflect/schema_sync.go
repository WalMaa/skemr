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
	"github.com/walmaa/skemr-api/internal/domain/databases"
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

type entitySyncInput struct {
	Name        string
	EntityType  sqlc.DatabaseEntityType
	ParentID    *uuid.UUID
	Fingerprint string
	Attributes  []byte
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
	err = s.SyncSchema(c, databases.ToDomainDatabase(database))
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

	// Get all schemaRefs in the database
	schemaRefs, err := connector.GetSchemas(c, conn)
	if err != nil {
		return fmt.Errorf("error getting schemaRefs: %w", err)
	}

	// For each schema, get tables, indexed and columns
	for _, schema := range schemaRefs {
		schema, err := s.updateNamespace(c, schema, database)
		if err != nil {
			return err
		}
		// Add schema ID to current entity ids
		slog.Debug("Adding schema to current entity ids", "schemaName", schema.Name, "schemaId", schema.ID)
		currentEntityIds = append(currentEntityIds, schema.ID)

		// Indexes
		indexes, err := connector.GetIndexesInSchema(c, conn, schema.Name)
		if err != nil {
			return fmt.Errorf("error getting indexes in schema %q: %w", schema.Name, err)
		}

		for _, indexRef := range indexes {
			index, err := s.updateIndex(c, indexRef, database, schema.ID)
			if err != nil {
				return fmt.Errorf("error updating indexes: %w", err)
			}
			// Add index ID to current entity ids
			slog.Debug("Adding index to current entity ids", "indexName", index.Name, "indexId", index.ID)
			currentEntityIds = append(currentEntityIds, index.ID)
		}

		// Tables
		tables, err := connector.GetTablesInSchema(c, conn, schema.Name)
		if err != nil {
			return fmt.Errorf("error getting tables in schema %q: %w", schema.Name, err)
		}

		for _, tableRef := range tables {
			table, err := s.updateTable(c, tableRef, database, schema.ID)

			if err != nil {
				return fmt.Errorf("error updating tables: %w", err)
			}
			// Add table ID to current entity ids
			slog.Debug("Adding table to current entity ids", "tableName", table.Name, "tableId", table.ID)
			currentEntityIds = append(currentEntityIds, table.ID)

					// Constraints
		constraints, err := connector.GetConstraintsInTable(c, conn, tableRef.Name)
		if err != nil {
			return fmt.Errorf("error getting constraints in schema %q: %w", schema.Name, err)
		}

		for _, constraintRef := range constraints {
			constraint, err := s.updateConstraint(c, constraintRef, database, table.ID)
			if err != nil {
				return fmt.Errorf("error updating constraints: %w", err)
			}
			// Add constraint ID to current entity ids
			slog.Debug("Adding constraint to current entity ids", "constraintName", constraint.Name, "constraintId", constraint.ID)
			currentEntityIds = append(currentEntityIds, constraint.ID)
		}

			// Columns
			columns, err := connector.ListColumnsInTable(c, conn, tableRef)
			if err != nil {
				return fmt.Errorf("error getting columns in table %q.%q: %w", schema.Name, tableRef.Name, err)
			}
			for _, column := range columns {
				column, err := s.updateColumn(c, column, database, table.ID)
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
			err := s.markEntityAsDeleted(c, database.ID, entity.ID)
			if err != nil {
				slog.Error("Error marking entity as deleted", "entityId", entity.ID, "error", err)
				// Do not return error as we want to continue marking other entities as deleted
			}
		}
	}

	return nil
}

func (s *SchemaSyncService) updateConstraint(c context.Context, constraintRef ConstraintRef, database models.Database, tableId uuid.UUID) (sqlc.DatabaseEntity, error) {
	return s.updateEntity(c, database, entitySyncInput{
		Name:        constraintRef.ConstraintName,
		EntityType:  sqlc.DatabaseEntityTypeConstraint,
		ParentID:    &tableId,
		Fingerprint: GenerateConstraintFingerprint(constraintRef, tableId),
	})
}

func (s *SchemaSyncService) updateIndex(c context.Context, indexRef IndexRef, database models.Database, schemaId uuid.UUID) (sqlc.DatabaseEntity, error) {
	return s.updateEntity(c, database, entitySyncInput{
		Name:        indexRef.IndexName,
		EntityType:  sqlc.DatabaseEntityTypeIndex,
		ParentID:    &schemaId,
		Fingerprint: GenerateIndexFingerprint(indexRef, schemaId),
	})
}

func (s *SchemaSyncService) updateNamespace(c context.Context, schemaRef SchemaRef, database models.Database) (sqlc.DatabaseEntity, error) {
	return s.updateEntity(c, database, entitySyncInput{
		Name:        schemaRef.Name,
		EntityType:  sqlc.DatabaseEntityTypeSchema,
		ParentID:    nil, // Schemas do not have a parent entity
		Fingerprint: GenerateNamespaceFingerprint(schemaRef),
	})

}

func (s *SchemaSyncService) updateTable(c context.Context, tableRef TableRef, database models.Database, schemaId uuid.UUID) (sqlc.DatabaseEntity, error) {

	return s.updateEntity(c, database, entitySyncInput{
		Name:        tableRef.Name,
		EntityType:  sqlc.DatabaseEntityTypeTable,
		ParentID:    &schemaId,
		Fingerprint: GenerateTableFingerprint(tableRef),
	})
}

// updateColumn syncs a database column to a database entity.
func (s *SchemaSyncService) updateColumn(c context.Context, columnRef ColumnRef, database models.Database, tableId uuid.UUID) (sqlc.DatabaseEntity, error) {

	columnAttributes := ColumnAttributes{
		DataType:  columnRef.DataType,
		Default:   columnRef.Default,
		Nullable:  columnRef.Nullable,
		Updatable: columnRef.Updatable,
	}
	attributesJson, err := json.Marshal(columnAttributes)
	if err != nil {
		slog.Error("error marshalling column attributes", "error", err)
		return sqlc.DatabaseEntity{}, err
	}

	return s.updateEntity(c, database, entitySyncInput{
		Name:        columnRef.Name,
		EntityType:  sqlc.DatabaseEntityTypeColumn,
		ParentID:    &tableId,
		Fingerprint: GenerateColumnFingerprint(columnRef, tableId),
		Attributes:  attributesJson,
	})
}

// updateEntity checks if an entity with the given name exists for the database.
// If it does not exist, it creates a new entity and if it has been renamed, it updates the name of the existing entity.
func (s *SchemaSyncService) updateEntity(c context.Context, database models.Database, input entitySyncInput) (sqlc.DatabaseEntity, error) {
	// Check if the entity exists by name and parent ID
	args := sqlc.GetDatabaseEntityByDatabaseIdAndTypeAndParentAndNameParams{
		DatabaseID: database.ID,
		ParentID:   input.ParentID,
		EntityType: sqlc.DatabaseEntityType(input.EntityType),
		Name:       input.Name,
	}
	entity, err := s.db.GetDatabaseEntityByDatabaseIdAndTypeAndParentAndName(c, args)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("error getting entity", "error", err.Error())
		return sqlc.DatabaseEntity{}, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		// If not found by name, check by fingerprint to see if it is the same entity with an updated name.

		entity, err = s.db.GetDatabaseEntityByFingerprint(c, sqlc.GetDatabaseEntityByFingerprintParams{
			DatabaseID:  database.ID,
			Fingerprint: input.Fingerprint,
		})

		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			slog.Error("error getting entity by fingerprint", "error", err.Error())
			return sqlc.DatabaseEntity{}, err
		}

		// If found by fingerprint, update the name to the new name.
		if !errors.Is(err, pgx.ErrNoRows) {
			slog.Debug("Entity found by fingerprint, updating name", "oldName", entity.Name, "newName", input.Name)

			entity, err = s.db.UpdateDatabaseEntityName(c, sqlc.UpdateDatabaseEntityNameParams{
				ID:   entity.ID,
				Name: input.Name,
			})
			if err != nil {
				slog.Error("error updating entity name", "error", err.Error())
				return sqlc.DatabaseEntity{}, err
			}

			s.markDatabaseChange(c, database.ID, entity.ID, models.MigrationStatementActionUpdate)
			slog.Info("Entity renamed", "oldName", entity.Name, "newName", input.Name)
			entity.Name = input.Name
			return entity, nil
		}

		// If that entity does not exist yet, save it
		args := sqlc.CreateDatabaseEntityParams{
			ProjectID:   database.ProjectID,
			EntityType:  sqlc.DatabaseEntityType(input.EntityType),
			ParentID:    input.ParentID,
			DatabaseID:  database.ID,
			Name:        input.Name,
			Fingerprint: input.Fingerprint,
		}
		entity, err := s.db.CreateDatabaseEntity(c, args)
		if err != nil {
			slog.Error("error creating entity", "error", err)
			return sqlc.DatabaseEntity{}, err
		}

		s.markDatabaseChange(c, database.ID, entity.ID, models.MigrationStatementActionCreate)
		slog.Info("Entity created", "name", entity.Name)
		return entity, err
	}

	// If entity does exist, update if needed. Currently no-op
	slog.Info("entity exists", "entity", entity.Name)

	return entity, err
}

// markEntityAsDeleted sets the status of the entity to "deleted".
// This is used for entities that were not found in the new schema during sync.
// We want to keep these entities in the database to show how the schema has changed.
func (s *SchemaSyncService) markEntityAsDeleted(c context.Context, databaseId uuid.UUID, entityId uuid.UUID) error {
	slog.Debug("Marking entity as deleted", "entityId", entityId)
	err := s.db.UpdateDatabaseEntityAsDeleted(c, entityId)
	if err != nil {
		slog.Error("error marking entity as deleted", "error", err)
		return err
	}

	s.markDatabaseChange(c, databaseId, entityId, models.MigrationStatementActionDelete)

	return nil
}

func (s *SchemaSyncService) markDatabaseChange(c context.Context, databaseId uuid.UUID, entityId uuid.UUID, action models.MigrationStatementAction) {
	slog.Debug("Marking database change", "databaseId", databaseId, "entityId", entityId)
	_, err := s.db.CreateDatabaseChange(c, sqlc.CreateDatabaseChangeParams{
		DatabaseID: databaseId,
		EntityID:   entityId,
		Action:     sqlc.MigrationStatementAction(action),
	})

	if err != nil {
		slog.Error("error creating database change", "error", err)
	}
}
