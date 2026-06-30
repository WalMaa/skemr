package ai

import (
	"context"

	"github.com/walmaa/skemr-api/internal/dbreflect"
	"github.com/walmaa/skemr-common/models"
)


type DatabaseReader interface {
	ReadDatabase(ctx context.Context, database models.Database) (string error)
}


type PostgresDatabaseReader struct {
	postgresConnector dbreflect.PostgresConnector
}

func NewPostgresDatabaseReader(database models.Database) *PostgresDatabaseReader {
	return &PostgresDatabaseReader{
		postgresConnector: dbreflect.PostgresConnector{Database: database},
	}
}

func (r *PostgresDatabaseReader) ReadDatabase(ctx context.Context, database models.Database) (string, error) {
	conn, err := r.postgresConnector.Connect(ctx)
	if err != nil {
		return "", err
	}
	defer r.postgresConnector.Disconnect(ctx, conn)

	conn.