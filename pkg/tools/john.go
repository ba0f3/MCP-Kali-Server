package tools

import (
	"fmt"

	"github.com/ba0f3/MCP-Kali-Server/pkg/executor"
)

// JohnParams represents parameters for John the Ripper
type JohnParams struct {
	HashFile       string `json:"hash_file"`
	Wordlist       string `json:"wordlist"`
	Format         string `json:"format"`
	AdditionalArgs string `json:"additional_args"`
}

// JohnCrack executes John the Ripper with the provided parameters
func JohnCrack(params JohnParams) (*ToolResult, error) {
	if params.HashFile == "" {
		return nil, fmt.Errorf("hash_file parameter is required")
	}

	// Default values
	if params.Wordlist == "" {
		params.Wordlist = "/usr/share/wordlists/rockyou.txt"
	}

	command := "john"

	if params.Format != "" {
		command += fmt.Sprintf(" --format=%s", params.Format)
	}

	if params.Wordlist != "" {
		command += fmt.Sprintf(" --wordlist=%s", params.Wordlist)
	}

	if params.AdditionalArgs != "" {
		command += fmt.Sprintf(" %s", params.AdditionalArgs)
	}

	command += fmt.Sprintf(" %s", params.HashFile)

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
