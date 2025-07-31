package tools

import (
	"fmt"

	"github.com/ba0f3/MCP-Kali-Server/pkg/executor"
)

// Sublist3rParams represents parameters for Sublist3r subdomain enumeration
type Sublist3rParams struct {
	Domain         string `json:"domain"`
	BruteForce     bool   `json:"bruteforce"`
	Ports          string `json:"ports"`
	Threads        int    `json:"threads"`
	Engines        string `json:"engines"`
	Verbose        bool   `json:"verbose"`
	AdditionalArgs string `json:"additional_args"`
}

// Sublist3rScan executes Sublist3r for subdomain enumeration
func Sublist3rScan(params Sublist3rParams) (*ToolResult, error) {
	if params.Domain == "" {
		return nil, fmt.Errorf("domain parameter is required")
	}

	// Build command
	command := "sublist3r"
	
	// Add domain
	command += fmt.Sprintf(" -d %s", params.Domain)

	// Add optional parameters
	if params.BruteForce {
		command += " -b"
	}

	if params.Ports != "" {
		command += fmt.Sprintf(" -p %s", params.Ports)
	}

	if params.Threads > 0 {
		command += fmt.Sprintf(" -t %d", params.Threads)
	} else {
		// Default threads
		command += " -t 10"
	}

	if params.Engines != "" {
		command += fmt.Sprintf(" -e %s", params.Engines)
	}

	if params.Verbose {
		command += " -v"
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
