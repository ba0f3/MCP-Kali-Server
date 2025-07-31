package tools

import (
	"fmt"

	"github.com/ba0f3/MCP-Kali-Server/pkg/executor"
)

// ToolResult represents the result of a tool execution
type ToolResult struct {
	Stdout         string `json:"stdout"`
	Stderr         string `json:"stderr"`
	Success        bool   `json:"success"`
	Error          string `json:"error,omitempty"`
	ReturnCode     int    `json:"return_code"`
	TimedOut       bool   `json:"timed_out"`
	PartialResults bool   `json:"partial_results"`
}

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

	// Validate mode
	validModes := map[string]bool{"dir": true, "dns": true, "fuzz": true, "vhost": true}
	if !validModes[params.Mode] {
		return nil, fmt.Errorf("invalid mode: %s. Must be one of: dir, dns, fuzz, vhost", params.Mode)
	}

	command := fmt.Sprintf("gobuster %s -u %s -w %s", params.Mode, params.URL, params.Wordlist)
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

// DirbParams represents parameters for Dirb scan
type DirbParams struct {
	URL            string `json:"url"`
	Wordlist       string `json:"wordlist"`
	AdditionalArgs string `json:"additional_args"`
}

// DirbScan executes Dirb with the provided parameters
func DirbScan(params DirbParams) (*ToolResult, error) {
	if params.URL == "" {
		return nil, fmt.Errorf("URL parameter is required")
	}

	// Default values
	if params.Wordlist == "" {
		params.Wordlist = "/usr/share/wordlists/dirb/common.txt"
	}

	command := fmt.Sprintf("dirb %s %s", params.URL, params.Wordlist)
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

// NiktoParams represents parameters for Nikto scan
type NiktoParams struct {
	Target         string `json:"target"`
	AdditionalArgs string `json:"additional_args"`
}

// NiktoScan executes Nikto with the provided parameters
func NiktoScan(params NiktoParams) (*ToolResult, error) {
	if params.Target == "" {
		return nil, fmt.Errorf("target parameter is required")
	}

	command := fmt.Sprintf("nikto -h %s", params.Target)
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

	command := fmt.Sprintf("hydra -t 4")

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

	command := fmt.Sprintf("john")

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

// WpscanParams represents parameters for WPScan
type WpscanParams struct {
	URL            string `json:"url"`
	AdditionalArgs string `json:"additional_args"`
}

// WpscanAnalyze executes WPScan with the provided parameters
func WpscanAnalyze(params WpscanParams) (*ToolResult, error) {
	if params.URL == "" {
		return nil, fmt.Errorf("URL parameter is required")
	}

	command := fmt.Sprintf("wpscan --url %s", params.URL)
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
