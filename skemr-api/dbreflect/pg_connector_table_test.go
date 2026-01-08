package dbreflect

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPostgresTablesInSchema(t *testing.T) {

	ctx, dbConn, conn := newTestPGConn(t)

	tables, err := dbConn.GetTablesInSchema(ctx, conn, "public")
	require.NoError(t, err)
	require.NotNil(t, tables)
	require.IsType(t, []TableRef{}, tables)

	err = conn.Close(ctx)
	require.NoError(t, err)
}

func TestGetPostgresColumn(t *testing.T) {
	ctx, dbConn, conn := newTestPGConn(t)

	columns, err := dbConn.ListColumnsInTable(ctx, conn, TableRef{Schema: "public", Name: "customers"})
	require.NoError(t, err)
	require.NotNil(t, columns)
	require.IsType(t, []ColumnRef{}, columns)

	err = conn.Close(ctx)
	require.NoError(t, err)
}

func TestGetPostgresSchema(t *testing.T) {
	ctx, dbConn, conn := newTestPGConn(t)

	schemas, err := dbConn.GetSchemas(ctx, conn)
	require.NoError(t, err)
	require.NotNil(t, schemas)
	require.IsType(t, []string{}, schemas)

	err = conn.Close(ctx)
	require.NoError(t, err)
}
