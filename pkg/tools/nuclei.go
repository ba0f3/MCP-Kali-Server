package tools

import (
	"fmt"

	"github.com/ba0f3/MCP-Kali-Server/pkg/executor"
	"github.com/ba0f3/MCP-Kali-Server/pkg/security"
)

// NucleiParams represents parameters for Nuclei scan
type NucleiParams struct {
	Target         string `json:"target"`
	Templates      string `json:"templates"`
	Severity       string `json:"severity"`
	Tags           string `json:"tags"`
	AdditionalArgs string `json:"additional_args"`
}

// NucleiScan executes Nuclei with the provided parameters
func NucleiScan(params NucleiParams) (*ToolResult, error) {
	if params.Target == "" {
		return nil, fmt.Errorf("target parameter is required")
	}

	// Sanitize arguments
	params.Target = security.SanitizeInput(params.Target)
	params.Templates = security.SanitizeInput(params.Templates)
	params.Tags = security.SanitizeInput(params.Tags)
	params.AdditionalArgs = security.SanitizeInput(params.AdditionalArgs)

	// Build command
	command := "nuclei"

	// Add target
	command += fmt.Sprintf(" -u %s", params.Target)

	// Add templates if specified
	if params.Templates != "" {
		command += fmt.Sprintf(" -t %s", params.Templates)
	}

	// Add severity filter if specified
	if params.Severity != "" {
		// Validate severity levels
		validSeverities := map[string]bool{
			"info": true, "low": true, "medium": true,
			"high": true, "critical": true,
		}
		if validSeverities[params.Severity] {
			command += fmt.Sprintf(" -s %s", params.Severity)
		} else {
			return nil, fmt.Errorf("invalid severity level: %s. Must be one of: info, low, medium, high, critical", params.Severity)
		}
	}

	// Add tags filter if specified
	if params.Tags != "" {
		command += fmt.Sprintf(" -tags %s", params.Tags)
	}

	// Add default flags for better output
	command += " -silent -no-interactsh"

	// Add any additional arguments
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
