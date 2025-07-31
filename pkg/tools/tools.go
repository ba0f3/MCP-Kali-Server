package tools

import (
	"fmt"

	"github.com/ba0f3/MCP-Kali-Server/pkg/executor"
	"github.com/ba0f3/MCP-Kali-Server/pkg/security"
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

	// if URL not starts with https://, replace -u with -do
	if params.Mode == "dns" {
		params.URL = "-do " + params.URL
	} else {
		params.URL = "-u " + params.URL
	}
	command := fmt.Sprintf("gobuster %s %s -w %s", params.Mode, params.URL, params.Wordlist)
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

// PingParams represents parameters for ping
type PingParams struct {
	Target         string `json:"target"`
	Count          int    `json:"count"`
	Timeout        int    `json:"timeout"`
	PacketSize     int    `json:"packet_size"`
	AdditionalArgs string `json:"additional_args"`
}

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
