//go:build !darwin
// +build !darwin

package service

import "fmt"

// installDarwinService is a stub for non-Darwin platforms
func installDarwinService(config ServiceConfig) error {
	return fmt.Errorf("macOS service installation is not supported on this platform")
}

// uninstallDarwinService is a stub for non-Darwin platforms
func uninstallDarwinService(name string) error {
	return fmt.Errorf("macOS service uninstallation is not supported on this platform")
}
