package ai

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

type ToolCall struct {
	ID        string          `json:"id"`
	CallID    string          `json:"callId"`
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

type Completion struct {
	Text      string     `json:"text"`
	ToolCalls []ToolCall `json:"toolCalls,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Model interface {
	Complete(ctx context.Context, msgs []Message, tools []ToolSpec, actor Actor) (Completion, error)
}

type OpenAIClient struct {
	client       *openai.Client
	toolRegistry *ToolRegistry
}

type Actor struct {
	UserID    uuid.UUID
	ProjectID uuid.UUID
}

func NewOpenAIClient(toolRegistry *ToolRegistry) *OpenAIClient {
	client := openai.NewClient()
	return &OpenAIClient{client: &client, toolRegistry: toolRegistry}
}

func (c *OpenAIClient) Complete(ctx context.Context, msgs []Message, tools []ToolSpec, actor Actor) (Completion, error) {

	params := responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String("What is the current working directory and list the files in it?"),
		},
		Model: openai.ChatModelGPT5Nano,
		Tools: c.toolRegistry.toToolUnionParams(),
	}

	response, err := c.prompt(ctx, params)

	if err != nil {
		slog.Error("Error generating response", "err", err)
		return Completion{}, err
	}

	var outputs []responses.ResponseInputItemUnionParam
	// Tool call handling
	for _, item := range response.Output {
		if item.Type != "function_call" {
			continue
		}
		toolCall := item.AsFunctionCall()
		toolCallResult, err := c.toolRegistry.Run(ctx, toolCall.Name, nil, actor)

		if err != nil {
			slog.Error("Error running tool", "toolName", toolCall.Name, "err", err)
		}
		slog.Debug("toolCallResult", "result", toolCallResult)

		responseInput := responses.ResponseInputItemUnionParam{
			OfFunctionCallOutput: &responses.ResponseInputItemFunctionCallOutputParam{
				CallID: toolCall.CallID,
				Output: responses.ResponseInputItemFunctionCallOutputOutputUnionParam{
					OfString: openai.String(toolCallResult),
				},
			},
		}

		outputs = append(outputs, responseInput)

	}
	toolCallResponseParams := responses.ResponseNewParams{
		Model:              openai.ChatModelGPT5Nano,
		PreviousResponseID: openai.String(response.ID),
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: outputs,
		},
	}

	// Continue conversation with tool call result
	toolCallResponse, err := c.prompt(ctx, toolCallResponseParams)

	if err != nil {
		slog.Error("Error generating response after tool call", "err", err)
		return Completion{}, err
	}

	return Completion{
		Text: toolCallResponse.OutputText(),
	}, nil
}

// prompt sends a prompt to the OpenAI API and returns the response.
// Used as a wrapper for logging
func (c *OpenAIClient) prompt(ctx context.Context, params responses.ResponseNewParams) (*responses.Response, error) {
	response, err := c.client.Responses.New(ctx, params)
	slog.Debug("response", "inputTokens", response.Usage.InputTokens, "outputTokens", response.Usage.OutputTokens, "totalTokens", response.Usage.TotalTokens)
	return response, err
}
