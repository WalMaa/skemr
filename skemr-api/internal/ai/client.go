package ai

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

type ToolCall struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Input json.RawMessage `json:"input"`
}

type Completion struct {
	Text string `json:"text"`
	ToolCalls []ToolCall `json:"toolCalls"`
}

type Message struct {
	Role string `json:"role"`
	Content string `json:"content"`
}

type ToolSpec struct {
	Name string `json:"name"`
	Description string `json:"description"`
}

type Model interface {
	Complete (ctx context.Context, msgs []Message, tools []ToolSpec) (Completion, error)
}

type Tool interface {
	Spec() ToolSpec
	Run(ctx context.Context, input json.RawMessage) (string, error)
}

type OpenAIClient struct {
	client *openai.Client
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	client := openai.NewClient()
	return &OpenAIClient{client: &client}
}



func (c *OpenAIClient) Complete(ctx context.Context, msgs []Message, tools []ToolSpec) (Completion, error) {

	params := responses.ResponseNewParams{
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String("Hello"),
			},
			Model: openai.ChatModelGPT5Nano,
			Tools: Tools,
		}

	response, err := c.client.Responses.New(ctx, params)

	if err != nil {
		slog.Error("Error generating response", "err", err)
		return Completion{}, err
	}

	for _, item := range response.Output {
		if item.Type == "function_call" {
			toolCall := item.AsFunctionCall()
			toolCallResult := ""
			switch toolCall.Name {
			case "get_cwd":
				cwd, _ := getCwd(ctx)
				slog.Info("Tool call result", "toolName", toolCall.Name, "result", cwd)
				toolCallResult = cwd
			case "ls":
				lsResult, _ := ls(ctx)
				slog.Info("Tool call result", "toolName", toolCall.Name, "result", lsResult)
				toolCallResult = lsResult
			default:
				slog.Warn("Unknown tool call", "toolName", toolCall.Name)
			}

			// Continue conversation with tool call result
			toolCallResponseParams := responses.ResponseNewParams{
				Model: openai.ChatModelGPT5Nano,
				PreviousResponseID: openai.String(response.ID),
				Input: responses.ResponseNewParamsInputUnion{
					OfInputItemList: []responses.ResponseInputItemUnionParam{{
						OfFunctionCallOutput: &responses.ResponseInputItemFunctionCallOutputParam{
							CallID: toolCall.ID,
							Output: responses.ResponseInputItemFunctionCallOutputOutputUnionParam{
								OfString: openai.String(toolCallResult),
							},
						},
					}},
				},
			}
			toolCallResponse, err := c.client.Responses.New(ctx, toolCallResponseParams)
			if err != nil {
				slog.Error("Error generating response after tool call", "err", err)
				return Completion{}, err
			}

			return Completion{
				Text: toolCallResponse.Output[0].AsMessage().JSON.Content.Raw(),
				ToolCalls: []ToolCall{{
					ID: toolCall.ID,
					Name: toolCall.Name,
					Input: nil,
				}},
			}, nil
		}
	}

	return Completion{
		Text: response.Output[0].AsMessage().JSON.Content.Raw(),
	}, nil
}
