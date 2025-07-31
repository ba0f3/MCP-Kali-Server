package tools

import (
	"fmt"

	"github.com/ba0f3/MCP-Kali-Server/pkg/executor"
)

// NmapParams represents parameters for Nmap scan
type NmapParams struct {
	Target         string `json:"target"`
	ScanType       string `json:"scan_type"`
	Ports          string `json:"ports"`
	AdditionalArgs string `json:"additional_args"`
}

// NmapScan executes an Nmap scan with the provided parameters
func NmapScan(params NmapParams) (*ToolResult, error) {
	if params.Target == "" {
		return nil, fmt.Errorf("target parameter is required")
	}

	// Default values
	if params.ScanType == "" {
		params.ScanType = "-sCV"
	}
	if params.AdditionalArgs == "" {
		params.AdditionalArgs = "-T4 -Pn"
	}

	command := fmt.Sprintf("nmap %s", params.ScanType)
	if params.Ports != "" {
		command += fmt.Sprintf(" -p %s", params.Ports)
	}
	if params.AdditionalArgs != "" {
		command += fmt.Sprintf(" %s", params.AdditionalArgs)
	}
	command += fmt.Sprintf(" %s", params.Target)

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
