package ai

import (
	"context"
	"encoding/json"

	"github.com/walmaa/skemr-common/models"
)

type DatabaseToolService interface {
	GetDatabases(ctx context.Context, actor Actor) ([]models.Database, error)
}

type DatabaseTool struct {
	service DatabaseToolService
}

func NewDatabaseTool(service DatabaseToolService) *DatabaseTool {
	return &DatabaseTool{service: service}
}

func (t *DatabaseTool) Spec() ToolSpec {
	return ToolSpec{
		Name:        "get_databases",
		Description: "Get the databases for a given project",
		Parameters:  nil,
		Strict:      true,
	}
}

func (t *DatabaseTool) Run(ctx context.Context, input json.RawMessage, actor Actor) (string, error) {
	databases, err := t.service.GetDatabases(ctx, actor)
	if err != nil {
		return toolError("search_failed", err.Error()), nil
	}

	return toolJSON(databases)
}
