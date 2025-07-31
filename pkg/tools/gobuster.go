package tools

import (
	"fmt"

	"github.com/ba0f3/MCP-Kali-Server/pkg/executor"
)

// GobusterParams represents parameters for Gobuster scan
type GobusterParams struct {
	URL            string `json:"url"`
	Mode           string `json:"mode"`
	Wordlist       string `json:"wordlist"`
	AdditionalArgs string `json:"additional_args"`
}

// GobusterScan executes Gobuster with the provided parameters
func GobusterScan(params GobusterParams) (*ToolResult, error) {
	if params.URL == "" {
		return nil, fmt.Errorf("URL parameter is required")
	}

	// Default values
	if params.Mode == "" {
		params.Mode = "dir"
	}
	if params.Wordlist == "" {
		params.Wordlist = "/usr/share/wordlists/dirb/common.txt"
	}

	var url string
	if params.Mode == "dns" {
		url = "-do " + params.URL
	} else {
		url = "-u " + params.URL
	}

	// Validate mode
	validModes := map[string]bool{"dir": true, "dns": true, "fuzz": true, "vhost": true}
	if !validModes[params.Mode] {
		return nil, fmt.Errorf("invalid mode: %s. Must be one of: dir, dns, fuzz, vhost", params.Mode)
	}

	command := fmt.Sprintf("gobuster %s %s -w %s", params.Mode, url, params.Wordlist)
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
