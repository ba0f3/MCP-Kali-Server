//go:build windows
// +build windows

package service

import (
	"fmt"
	"os/exec"
	"strings"
)

// installWindowsService installs the service on Windows using sc.exe
func installWindowsService(config ServiceConfig) error {
	// Build the binpath with arguments
	binPath := fmt.Sprintf("\"%s\"", config.Executable)
	if len(config.Args) > 0 {
		binPath += " " + strings.Join(config.Args, " ")
	}

	// Create the service
	cmd := exec.Command("sc.exe", "create", config.Name,
		"binPath=", binPath,
		"DisplayName=", config.DisplayName,
		"start=", "auto")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create service: %v, output: %s", err, string(output))
	}

	// Set service description
	cmd = exec.Command("sc.exe", "description", config.Name, config.Description)
	output, err = cmd.CombinedOutput()
	if err != nil {
		// Non-critical error, just log it
		fmt.Printf("Warning: failed to set service description: %v\n", err)
	}

	// Configure failure actions (restart on failure)
	cmd = exec.Command("sc.exe", "failure", config.Name,
		"reset=", "86400",
		"actions=", "restart/60000/restart/60000/restart/60000")
	output, err = cmd.CombinedOutput()
	if err != nil {
		// Non-critical error, just log it
		fmt.Printf("Warning: failed to set failure actions: %v\n", err)
	}

	return nil
}

// uninstallWindowsService removes the service on Windows
func uninstallWindowsService(name string) error {
	// Stop the service first
	cmd := exec.Command("sc.exe", "stop", name)
	cmd.Run() // Ignore error if service is not running

	// Delete the service
	cmd = exec.Command("sc.exe", "delete", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete service: %v, output: %s", err, string(output))
	}

	return nil
}
