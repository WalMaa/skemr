package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

type Tool interface {
	Spec() ToolSpec
	Run(ctx context.Context, input json.RawMessage, actor Actor) (string, error)
}

type ToolSpec struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
	Strict      bool            `json:"strict"`
}

type ToolRegistry struct {
	tools map[string]Tool
}

func NewToolRegistry(tools ...Tool) *ToolRegistry {
	toolMap := make(map[string]Tool)

	for _, tool := range toolMap {
		toolMap[tool.Spec().Name] = tool
	}
	return &ToolRegistry{
		tools: toolMap,
	}
}

func (r *ToolRegistry) toToolUnionParams() []responses.ToolUnionParam {
	var toolUnionParams []responses.ToolUnionParam
	for _, tool := range r.tools {
		spec := tool.Spec()
		toolUnionParams = append(toolUnionParams, responses.ToolUnionParam{
			OfFunction: &responses.FunctionToolParam{
				Name:        spec.Name,
				Description: openai.String(spec.Description),
			},
		})
	}
	return toolUnionParams
}

// Run executes the specified tool with the given input and actor context.
func (r *ToolRegistry) Run(ctx context.Context, toolName string, input json.RawMessage, actor Actor) (string, error) {
	tool, ok := r.tools[toolName]
	if !ok {
		return "", fmt.Errorf("unknown tool: %s", toolName)
	}
	return tool.Run(ctx, input, actor)
}
