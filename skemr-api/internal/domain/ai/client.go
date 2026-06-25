package ai

import (
	"context"
	"encoding/json"
	"fmt"
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
	Result    string          `json:"result,omitempty"`
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
	Complete(ctx context.Context, msgs []Message, actor Actor) (Completion, error)
}

type OpenAIClient struct {
	client       *openai.Client
	toolRegistry *ToolRegistry
}

type Actor struct {
	UserID    uuid.UUID
	ProjectID uuid.UUID
}

const maxToolTurns = 5
const defaultModel = openai.ChatModelGPT5Nano

const systemPrompt = `You are a helpful assistant that can answer questions about the current state of the system.
 You have access to a set of tools that can provide information about the system's databases and their contents. 
 Use these tools to answer questions accurately and efficiently.
 If you can not answer a question with the available tools, respond with "I don't know" instead of making up an answer.
 Do not make up any information.
 Do not claim functionality that does not exist within the provided tooling.`

func NewOpenAIClient(toolRegistry *ToolRegistry) *OpenAIClient {
	client := openai.NewClient()
	return &OpenAIClient{client: &client, toolRegistry: toolRegistry}
}

func (c *OpenAIClient) collectAndRunToolCalls(ctx context.Context, response *responses.Response, actor Actor) ([]ToolCall, error) {
	var toolCalls []ToolCall

	for _, item := range response.Output {
		if item.Type != "function_call" {
			continue
		}
		toolCall := item.AsFunctionCall()
		slog.Debug("Calling tool", "toolName", toolCall.Name)
		toolCallResult, err := c.toolRegistry.Run(ctx, toolCall.Name, json.RawMessage(toolCall.Arguments), actor)

		if err != nil {
			slog.Error("Error running tool", "toolName", toolCall.Name, "err", err)
			return nil, err
		}
		slog.Debug("toolCallResult", "result", toolCallResult)

		toolCalls = append(toolCalls, ToolCall{
			ID:        toolCall.ID,
			CallID:    toolCall.CallID,
			Name:      toolCall.Name,
			Arguments: nil,
			Result:    toolCallResult,
		})
	}

	return toolCalls, nil
}

func (c *OpenAIClient) Complete(ctx context.Context, msgs []Message, actor Actor) (Completion, error) {
	slog.Info("Completing with OpenAI", "actor", actor, "messages", msgs)

	if len(msgs) == 0 {
		return Completion{}, fmt.Errorf("no messages provided")
	}

	userPrompt := msgs[len(msgs)-1].Content

	initialParams := responses.ResponseNewParams{
		Instructions: openai.String(systemPrompt),
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(userPrompt),
		},
		Model: defaultModel,
		Tools: c.toolRegistry.toToolUnionParams(),
	}

	params := initialParams

	for range maxToolTurns {
		// On first loop iteration, use the initial params. On subsequent iterations, use the params from the previous response.
		// In this way, we can continue the conversation with the tool call results.
		response, err := c.prompt(ctx, params)

		if err != nil {
			slog.Error("Error generating response", "err", err)
			return Completion{}, err
		}

		// Tool call handling
		toolCalls, err := c.collectAndRunToolCalls(ctx, response, actor)

		if err != nil {
			slog.Error("Error collecting and running tool calls", "err", err)
			return Completion{}, err
		}

		var outputs []responses.ResponseInputItemUnionParam

		for _, toolCall := range toolCalls {

			responseInput := responses.ResponseInputItemUnionParam{
				OfFunctionCallOutput: &responses.ResponseInputItemFunctionCallOutputParam{
					CallID: toolCall.CallID,
					Output: responses.ResponseInputItemFunctionCallOutputOutputUnionParam{
						OfString: openai.String(toolCall.Result),
					},
				},
			}

			outputs = append(outputs, responseInput)

		}
		toolCallResponseParams := responses.ResponseNewParams{
			Instructions:       openai.String(systemPrompt),
			Model:              defaultModel,
			Tools:              c.toolRegistry.toToolUnionParams(),
			PreviousResponseID: openai.String(response.ID),
			Input: responses.ResponseNewParamsInputUnion{
				OfInputItemList: outputs,
			},
		}

		// If there are no tool calls, return the original response
		if len(outputs) == 0 {
			return Completion{
				Text: response.OutputText(),
			}, nil
		}

		// Update params for the next iteration
		params = toolCallResponseParams
	}

	return Completion{}, fmt.Errorf("Tool loop exceeded %d turns", maxToolTurns)

}

// prompt calls the OpenAI API and returns the response.
// Used as a wrapper for logging
func (c *OpenAIClient) prompt(ctx context.Context, params responses.ResponseNewParams) (*responses.Response, error) {
	response, err := c.client.Responses.New(ctx, params)

	if err == nil {
		slog.Debug("response", "inputTokens", response.Usage.InputTokens, "outputTokens", response.Usage.OutputTokens, "totalTokens", response.Usage.TotalTokens)
	}
	return response, err
}
