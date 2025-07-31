package service

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// ServiceConfig holds the configuration for the service
type ServiceConfig struct {
	Name        string
	DisplayName string
	Description string
	Executable  string
	Args        []string
	WorkingDir  string
}

// InstallService installs the application as a system service
func InstallService(config ServiceConfig) error {
	switch runtime.GOOS {
	case "windows":
		return installWindowsService(config)
	case "linux":
		return installLinuxService(config)
	case "darwin":
		return installDarwinService(config)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// UninstallService removes the installed service
func UninstallService(name string) error {
	switch runtime.GOOS {
	case "windows":
		return uninstallWindowsService(name)
	case "linux":
		return uninstallLinuxService(name)
	case "darwin":
		return uninstallDarwinService(name)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// GetDefaultConfig returns a default service configuration
func GetDefaultConfig() ServiceConfig {
	executable, _ := os.Executable()
	workingDir := filepath.Dir(executable)

	return ServiceConfig{
		Name:        "mcp-kali-server",
		DisplayName: "MCP Kali Server",
		Description: "Model Context Protocol server for Kali Linux security tools",
		Executable:  executable,
		WorkingDir:  workingDir,
	}
}
