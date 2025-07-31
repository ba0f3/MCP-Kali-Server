package tools

import (
	"fmt"

	"github.com/ba0f3/MCP-Kali-Server/pkg/executor"
)

// HydraParams represents parameters for Hydra attack
type HydraParams struct {
	Target         string `json:"target"`
	Service        string `json:"service"`
	Username       string `json:"username"`
	UsernameFile   string `json:"username_file"`
	Password       string `json:"password"`
	PasswordFile   string `json:"password_file"`
	AdditionalArgs string `json:"additional_args"`
}

// HydraAttack executes Hydra with the provided parameters
func HydraAttack(params HydraParams) (*ToolResult, error) {
	if params.Target == "" {
		return nil, fmt.Errorf("target parameter is required")
	}
	if params.Service == "" {
		return nil, fmt.Errorf("service parameter is required")
	}
	if params.Username == "" && params.UsernameFile == "" {
		return nil, fmt.Errorf("username or username_file parameter is required")
	}
	if params.Password == "" && params.PasswordFile == "" {
		return nil, fmt.Errorf("password or password_file parameter is required")
	}

	command := "hydra -t 4"

	if params.Username != "" {
		command += fmt.Sprintf(" -l %s", params.Username)
	} else {
		command += fmt.Sprintf(" -L %s", params.UsernameFile)
	}

	if params.Password != "" {
		command += fmt.Sprintf(" -p %s", params.Password)
	} else {
		command += fmt.Sprintf(" -P %s", params.PasswordFile)
	}

	if params.AdditionalArgs != "" {
		command += fmt.Sprintf(" %s", params.AdditionalArgs)
	}

	command += fmt.Sprintf(" %s %s", params.Target, params.Service)

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
