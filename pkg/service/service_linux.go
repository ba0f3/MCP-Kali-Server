//go:build linux
// +build linux

package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const systemdTemplate = `[Unit]
Description={{.Description}}
After=network.target

[Service]
Type=simple
User={{.User}}
WorkingDirectory={{.WorkingDir}}
ExecStart={{.Executable}}{{range .Args}} {{.}}{{end}}
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
`

// installLinuxService installs the service on Linux using systemd
func installLinuxService(config ServiceConfig) error {
	// Check if systemd is available
	if _, err := os.Stat("/usr/bin/systemctl"); os.IsNotExist(err) {
		return fmt.Errorf("systemd is not available on this system")
	}

	// Create service file
	serviceFile := filepath.Join("/etc/systemd/system", config.Name+".service")

	file, err := os.Create(serviceFile)
	if err != nil {
		return fmt.Errorf("failed to create service file: %v", err)
	}
	defer file.Close()

	// Prepare template data
	data := struct {
		Description string
		User        string
		WorkingDir  string
		Executable  string
		Args        []string
	}{
		Description: config.Description,
		User:        os.Getenv("USER"),
		WorkingDir:  config.WorkingDir,
		Executable:  config.Executable,
		Args:        config.Args,
	}

	// If running as root, use nobody user
	if data.User == "root" || data.User == "" {
		data.User = "nobody"
	}

	// Write service file
	tmpl, err := template.New("service").Parse(systemdTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to write service file: %v", err)
	}

	// Reload systemd
	cmd := exec.Command("systemctl", "daemon-reload")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to reload systemd: %v", err)
	}

	// Enable the service
	cmd = exec.Command("systemctl", "enable", config.Name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to enable service: %v", err)
	}

	return nil
}

// uninstallLinuxService removes the service on Linux
func uninstallLinuxService(name string) error {
	// Stop the service
	cmd := exec.Command("systemctl", "stop", name)
	cmd.Run() // Ignore error if service is not running

	// Disable the service
	cmd = exec.Command("systemctl", "disable", name)
	cmd.Run() // Ignore error if service is not enabled

	// Remove service file
	serviceFile := filepath.Join("/etc/systemd/system", name+".service")
	if err := os.Remove(serviceFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove service file: %v", err)
	}

	// Reload systemd
	cmd = exec.Command("systemctl", "daemon-reload")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to reload systemd: %v", err)
	}

	return nil
}
