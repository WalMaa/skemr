package dbreflect

import (
	"fmt"

	"github.com/google/uuid"
)

// GenerateColumnFingerprint generates a unique identifier for a column based on its properties.
// The principle is to create a stable identifier that remains consistent across column renames and db instance changes (backup restores).
// The format is column:{parent_table_id}:{ordinal_position}:{data_type}:{nullable}
func GenerateColumnFingerprint(columnRef ColumnRef, tableId uuid.UUID) string {
	return fmt.Sprintf("column:%s:%d:%s:%s", tableId.String(), columnRef.OrdinalPosition, columnRef.DataType, columnRef.Nullable)
}

// GenerateTableFingerprint generates a unique identifier for a table based on its properties.
// The principle is to create a stable identifier that remains consistent across table renames and db instance changes (backup restores).
// The format is table:{column_shape}:{primary_key}
func GenerateTableFingerprint(tableRef TableRef) string {
	return fmt.Sprintf("table:%s:%s", tableRef.ColumnShape, tableRef.PrimaryKey)
}

// GenerateNamespaceFingerprint generates a unique identifier for a schema based on its properties.
// The principle is to create a stable identifier that remains consistent across schema renames and db instance changes (backup restores).
// The format is schema:{database_id}:{schema_fingerprint}
func GenerateNamespaceFingerprint(schemaRef SchemaRef) string {
	return fmt.Sprintf("namespace:%s", schemaRef.Fingerprint)
}

// GenerateIndexFingerprint generates a unique identifier for an index based on its properties.
// The principle is to create a stable identifier that remains consistent across index renames and db instance changes (backup restores).
// The format is index:{schema_id}:{index_type}:{is_primary}:{is_unique}:{index_name}
func GenerateIndexFingerprint(indexRef IndexRef, schemaId uuid.UUID) string {
	return fmt.Sprintf("index:%s:%s:%t:%t:%s", schemaId.String(), indexRef.IndexType, indexRef.IsPrimary, indexRef.IsUnique, indexRef.IndexName)
}

// GenerateConstraintFingerprint generates a unique identifier for a constraint based on its properties.
// The principle is to create a stable identifier that remains consistent across constraint renames and db instance changes (backup restores).
// The format is constraint:{table_id}:{constraint_type}:{constraint_name}
func GenerateConstraintFingerprint(constraintRef ConstraintRef, tableId uuid.UUID) string {
	return fmt.Sprintf("constraint:%s:%s:%s", tableId.String(), constraintRef.ConstraintType, constraintRef.ConstraintName)
}