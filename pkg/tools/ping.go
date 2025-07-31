package tools

import (
	"fmt"

	"github.com/ba0f3/MCP-Kali-Server/pkg/executor"
	"github.com/ba0f3/MCP-Kali-Server/pkg/security"
)

// PingParams represents parameters for ping
type PingParams struct {
	Target         string `json:"target"`
	Count          int    `json:"count"`
	Timeout        int    `json:"timeout"`
	PacketSize     int    `json:"packet_size"`
	AdditionalArgs string `json:"additional_args"`
}

// Ping executes ping command with the provided parameters
func Ping(params PingParams) (*ToolResult, error) {
	if params.Target == "" {
		return nil, fmt.Errorf("target parameter is required")
	}

	// Sanitize target
	target, err := security.SanitizeTarget(params.Target)
	if err != nil {
		return nil, fmt.Errorf("invalid target: %v", err)
	}

	// Default values
	if params.Count <= 0 {
		params.Count = 4
	}
	if params.Timeout <= 0 {
		params.Timeout = 5
	}

	// Build command
	command := "ping"

	// Add count parameter
	command += fmt.Sprintf(" -c %d", params.Count)

	// Add timeout parameter
	command += fmt.Sprintf(" -W %d", params.Timeout)

	// Add packet size if specified
	if params.PacketSize > 0 {
		if params.PacketSize > 65507 {
			return nil, fmt.Errorf("packet size too large (max: 65507)")
		}
		command += fmt.Sprintf(" -s %d", params.PacketSize)
	}

	// Add target
	command += fmt.Sprintf(" %s", target)

	// Add any additional arguments (sanitized)
	if params.AdditionalArgs != "" {
		sanitizedArgs := security.SanitizeInput(params.AdditionalArgs)
		command += fmt.Sprintf(" %s", sanitizedArgs)
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
