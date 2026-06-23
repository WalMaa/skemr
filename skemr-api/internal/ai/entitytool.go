package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/walmaa/skemr-common/models"
)

type DatabaseEntityToolService interface {
	GetDatabaseEntities(ctx context.Context, databaseId uuid.UUID, actor Actor) ([]models.DatabaseEntity, error)
}

type DatabaseEntityTool struct {
	service DatabaseEntityToolService
}

func NewDatabaseEntityTool(service DatabaseEntityToolService) *DatabaseEntityTool {

	return &DatabaseEntityTool{service: service}

}

func (t *DatabaseEntityTool) Spec() ToolSpec {

	return ToolSpec{

		Name: "get_database_entities",

		Description: "Get the database entities for a given project and database",

		Parameters: json.RawMessage(`{"projectId": "string", "databaseId": "string"}`),

		Strict: true,
	}

}

func (t *DatabaseEntityTool) Run(ctx context.Context, input json.RawMessage, actor Actor) (string, error) {

	var params struct {
		DatabaseId uuid.UUID `json:"databaseId"`
	}

	if err := json.Unmarshal(input, &params); err != nil {

		return "", fmt.Errorf("invalid input parameters: %w", err)

	}

	entities, err := t.service.GetDatabaseEntities(ctx, params.DatabaseId, actor)

	if err != nil {
		return toolError("search_failed", err.Error()), nil
	}

	return toolJSON(entities)

}
