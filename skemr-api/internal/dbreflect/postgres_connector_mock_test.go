package dbreflect

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

type MockPostgresConnector struct {
	mock.Mock
}

func (m *MockPostgresConnector) Connect(ctx context.Context) (*pgx.Conn, error) {
	args := m.Called(ctx)
	conn := args.Get(0)
	if conn == nil {
		return nil, args.Error(1)
	}
	return conn.(*pgx.Conn), args.Error(1)
}

func (m *MockPostgresConnector) Disconnect(c context.Context, conn *pgx.Conn) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPostgresConnector) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPostgresConnector) TestConnection(c context.Context) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPostgresConnector) GetSchemas(c context.Context, conn *pgx.Conn) ([]SchemaRef, error) {
	args := m.Called()
	return args.Get(0).([]SchemaRef), args.Error(1)
}

func (m *MockPostgresConnector) GetTablesInSchema(c context.Context, conn *pgx.Conn, schema string) ([]TableRef, error) {
	args := m.Called()
	return args.Get(0).([]TableRef), args.Error(1)
}

func (m *MockPostgresConnector) ListColumnsInTable(c context.Context, conn *pgx.Conn, table TableRef) ([]ColumnRef, error) {
	args := m.Called(table)
	return args.Get(0).([]ColumnRef), args.Error(1)
}

func (m *MockPostgresConnector) getConnectionString() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
