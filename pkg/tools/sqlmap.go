package tools

import (
	"fmt"

	"github.com/ba0f3/MCP-Kali-Server/pkg/executor"
)

// SqlmapParams represents parameters for SQLmap scan
type SqlmapParams struct {
	URL            string `json:"url"`
	Data           string `json:"data"`
	AdditionalArgs string `json:"additional_args"`
}

// SqlmapScan executes SQLmap with the provided parameters
func SqlmapScan(params SqlmapParams) (*ToolResult, error) {
	if params.URL == "" {
		return nil, fmt.Errorf("URL parameter is required")
	}

	command := fmt.Sprintf("sqlmap -u %s --batch", params.URL)
	if params.Data != "" {
		command += fmt.Sprintf(" --data=\"%s\"", params.Data)
	}
	if params.AdditionalArgs != "" {
		command += fmt.Sprintf(" %s", params.AdditionalArgs)
	}

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
