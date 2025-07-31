//go:build !linux
// +build !linux

package service

import "fmt"

// installLinuxService is a stub for non-Linux platforms
func installLinuxService(config ServiceConfig) error {
	return fmt.Errorf("Linux service installation is not supported on this platform")
}

// uninstallLinuxService is a stub for non-Linux platforms
func uninstallLinuxService(name string) error {
	return fmt.Errorf("Linux service uninstallation is not supported on this platform")
}
