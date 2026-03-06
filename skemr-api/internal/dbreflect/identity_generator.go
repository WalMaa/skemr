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
// The format is table:{schema_id}:{column_shape}:{primary_key}
func GenerateTableFingerprint(tableRef TableRef, schemaId uuid.UUID) string {
	return fmt.Sprintf("table:%s:%s:%s", schemaId.String(), tableRef.ColumnShape, tableRef.PrimaryKey)
}
