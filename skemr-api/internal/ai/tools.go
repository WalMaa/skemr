package ai

import (
	"context"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)


var Tools = []responses.ToolUnionParam{
	{
		OfFunction: &cwdTool,
	},
	{
		OfFunction: &lsTool,
	},
}


var cwdTool = responses.FunctionToolParam{
	Name: "get_cwd",
	Description: openai.String("Get the current working directory"),
}

var lsTool = responses.FunctionToolParam{
	Name: "ls",
	Description: openai.String("List the files of the current directory"),
}

func getCwd(ctx context.Context) (string, error) {
	return "/home/user", nil
}

func ls(ctx context.Context) (string, error) {
	return "file1.txt\nfile2.txt\n", nil
}