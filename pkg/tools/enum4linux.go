package tools

import (
	"fmt"

	"github.com/ba0f3/MCP-Kali-Server/pkg/executor"
)

// Enum4linuxParams represents parameters for Enum4linux
type Enum4linuxParams struct {
	Target         string `json:"target"`
	AdditionalArgs string `json:"additional_args"`
}

// Enum4linuxScan executes Enum4linux with the provided parameters
func Enum4linuxScan(params Enum4linuxParams) (*ToolResult, error) {
	if params.Target == "" {
		return nil, fmt.Errorf("target parameter is required")
	}

	// Default values
	if params.AdditionalArgs == "" {
		params.AdditionalArgs = "-a"
	}

	command := fmt.Sprintf("enum4linux %s %s", params.AdditionalArgs, params.Target)

	result, err := executor.ExecuteCommand(command)
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
