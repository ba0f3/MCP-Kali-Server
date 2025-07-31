//go:build !windows
// +build !windows

package service

import "fmt"

// installWindowsService is a stub for non-Windows platforms
func installWindowsService(config ServiceConfig) error {
	return fmt.Errorf("Windows service installation is not supported on this platform")
}

// uninstallWindowsService is a stub for non-Windows platforms
func uninstallWindowsService(name string) error {
	return fmt.Errorf("Windows service uninstallation is not supported on this platform")
}
