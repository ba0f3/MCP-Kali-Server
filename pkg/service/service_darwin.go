//go:build darwin
// +build darwin

package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const launchdTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>{{.Name}}</string>
    <key>ProgramArguments</key>
    <array>
        <string>{{.Executable}}</string>
        {{range .Args}}<string>{{.}}</string>
        {{end}}
    </array>
    <key>WorkingDirectory</key>
    <string>{{.WorkingDir}}</string>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/var/log/{{.Name}}.log</string>
    <key>StandardErrorPath</key>
    <string>/var/log/{{.Name}}.error.log</string>
</dict>
</plist>
`

// installDarwinService installs the service on macOS using launchd
func installDarwinService(config ServiceConfig) error {
	// Create plist file
	plistFile := filepath.Join("/Library/LaunchDaemons", config.Name+".plist")

	file, err := os.Create(plistFile)
	if err != nil {
		return fmt.Errorf("failed to create plist file: %v", err)
	}
	defer file.Close()

	// Write plist file
	tmpl, err := template.New("plist").Parse(launchdTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	if err := tmpl.Execute(file, config); err != nil {
		return fmt.Errorf("failed to write plist file: %v", err)
	}

	// Set proper permissions
	if err := os.Chmod(plistFile, 0644); err != nil {
		return fmt.Errorf("failed to set permissions: %v", err)
	}

	// Load the service
	cmd := exec.Command("launchctl", "load", plistFile)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to load service: %v", err)
	}

	return nil
}

// uninstallDarwinService removes the service on macOS
func uninstallDarwinService(name string) error {
	plistFile := filepath.Join("/Library/LaunchDaemons", name+".plist")

	// Unload the service
	cmd := exec.Command("launchctl", "unload", plistFile)
	cmd.Run() // Ignore error if service is not loaded

	// Remove plist file
	if err := os.Remove(plistFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove plist file: %v", err)
	}

	return nil
}
