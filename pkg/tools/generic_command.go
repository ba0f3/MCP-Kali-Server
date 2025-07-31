package tools

import (
	"fmt"

	"github.com/ba0f3/MCP-Kali-Server/pkg/executor"
)

// GenericCommandParams represents parameters for generic command execution
type GenericCommandParams struct {
	Command string `json:"command"`
}

// ExecuteGenericCommand executes any command
func ExecuteGenericCommand(params GenericCommandParams) (*ToolResult, error) {
	if params.Command == "" {
		return nil, fmt.Errorf("command parameter is required")
	}

	result, err := executor.ExecuteCommand(params.Command)
	if err != nil {
		return nil, err
	}

	return &ToolResult{
		Stdout:         result.Stdout,
		Stderr:         result.Stderr,
		Success:        result.Success,
		ReturnCode:     result.ReturnCode,
		TimedOut:       result.TimedOut,
		PartialResults: result.PartialResults,
	}, nil
}
