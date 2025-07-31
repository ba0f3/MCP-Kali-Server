package security

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

// Command injection prevention patterns
var (
	// Dangerous shell metacharacters that could lead to command injection
	dangerousChars = regexp.MustCompile(`[;&|<>$` + "`" + `\\\n\r"'{}()\[\]*?~!]`)
	
	// Allow only alphanumeric, dots, hyphens, underscores, forward slashes, and colons
	safeChars = regexp.MustCompile(`^[a-zA-Z0-9\.\-_/:]+$`)
	
	// IP address pattern
	ipPattern = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	
	// Domain name pattern
	domainPattern = regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)*[a-zA-Z]{2,}$`)
	
	// Port range pattern
	portPattern = regexp.MustCompile(`^(\d{1,5}(-\d{1,5})?)(,\d{1,5}(-\d{1,5})?)*$`)
	
	// File path pattern (Unix/Windows)
	filePathPattern = regexp.MustCompile(`^[a-zA-Z0-9\.\-_/\\: ]+$`)
)

// SanitizeInput removes or escapes potentially dangerous characters from input
func SanitizeInput(input string) string {
	if input == "" {
		return ""
	}
	
	// Remove any null bytes
	input = strings.ReplaceAll(input, "\x00", "")
	
	// Escape shell metacharacters
	return dangerousChars.ReplaceAllStringFunc(input, func(match string) string {
		return "\\" + match
	})
}

// SanitizeURL validates and sanitizes URL input
func SanitizeURL(urlStr string) (string, error) {
	if urlStr == "" {
		return "", fmt.Errorf("empty URL")
	}
	
	// Parse the URL
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("invalid URL format: %v", err)
	}
	
	// Only allow http and https schemes
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", fmt.Errorf("invalid URL scheme: %s (only http/https allowed)", u.Scheme)
	}
	
	// Validate host
	if u.Host == "" {
		return "", fmt.Errorf("missing host in URL")
	}
	
	// Return the normalized URL
	return u.String(), nil
}

// SanitizeTarget validates and sanitizes target input (IP, domain, or URL)
func SanitizeTarget(target string) (string, error) {
	if target == "" {
		return "", fmt.Errorf("empty target")
	}
	
	// Remove any whitespace
	target = strings.TrimSpace(target)
	
	// Check if it's a URL
	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		return SanitizeURL(target)
	}
	
	// Check if it's an IP address
	if ipPattern.MatchString(target) {
		parts := strings.Split(target, ".")
		for _, part := range parts {
			if val := parseInt(part); val < 0 || val > 255 {
				return "", fmt.Errorf("invalid IP address")
			}
		}
		return target, nil
	}
	
	// Check if it's a domain name
	if domainPattern.MatchString(target) {
		return target, nil
	}
	
	// If it contains dangerous characters, reject it
	if dangerousChars.MatchString(target) {
		return "", fmt.Errorf("target contains invalid characters")
	}
	
	return target, nil
}

// SanitizeFilePath validates and sanitizes file path input
func SanitizeFilePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("empty file path")
	}
	
	// Clean the path
	cleanPath := filepath.Clean(path)
	
	// Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		return "", fmt.Errorf("path traversal detected")
	}
	
	// Ensure the path doesn't contain dangerous characters
	if !filePathPattern.MatchString(cleanPath) {
		return "", fmt.Errorf("file path contains invalid characters")
	}
	
	return cleanPath, nil
}

// SanitizePorts validates and sanitizes port specification
func SanitizePorts(ports string) (string, error) {
	if ports == "" {
		return "", nil // Empty ports is often valid (means default)
	}
	
	// Remove spaces
	ports = strings.ReplaceAll(ports, " ", "")
	
	// Validate port pattern
	if !portPattern.MatchString(ports) {
		return "", fmt.Errorf("invalid port specification")
	}
	
	// Validate individual ports
	portRanges := strings.Split(ports, ",")
	for _, portRange := range portRanges {
		if strings.Contains(portRange, "-") {
			parts := strings.Split(portRange, "-")
			if len(parts) != 2 {
				return "", fmt.Errorf("invalid port range")
			}
			start := parseInt(parts[0])
			end := parseInt(parts[1])
			if start < 1 || start > 65535 || end < 1 || end > 65535 || start > end {
				return "", fmt.Errorf("invalid port range values")
			}
		} else {
			port := parseInt(portRange)
			if port < 1 || port > 65535 {
				return "", fmt.Errorf("invalid port number: %d", port)
			}
		}
	}
	
	return ports, nil
}

// SanitizeCommand validates and sanitizes generic command arguments
// This is more restrictive and should be used carefully
func SanitizeCommand(cmd string) (string, error) {
	if cmd == "" {
		return "", fmt.Errorf("empty command")
	}
	
	// Check for common injection patterns
	dangerousPatterns := []string{
		"&&", "||", ";", "|", "&",
		"$(", "`", "${",
		">", "<", ">>", "<<",
		"\n", "\r",
	}
	
	for _, pattern := range dangerousPatterns {
		if strings.Contains(cmd, pattern) {
			return "", fmt.Errorf("command contains dangerous pattern: %s", pattern)
		}
	}
	
	return cmd, nil
}

// SanitizeAlphanumeric allows only alphanumeric characters, hyphens, and underscores
func SanitizeAlphanumeric(input string) (string, error) {
	if input == "" {
		return "", nil
	}
	
	alphanumPattern := regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)
	if !alphanumPattern.MatchString(input) {
		return "", fmt.Errorf("input contains non-alphanumeric characters")
	}
	
	return input, nil
}

// ValidateEnum validates that input is one of the allowed values
func ValidateEnum(input string, allowedValues []string) error {
	for _, allowed := range allowedValues {
		if input == allowed {
			return nil
		}
	}
	return fmt.Errorf("invalid value: %s (allowed: %v)", input, allowedValues)
}

// Helper function to parse integers safely
func parseInt(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}

// EscapeShellArg escapes a string for safe use as a shell argument
func EscapeShellArg(arg string) string {
	if arg == "" {
		return "''"
	}
	
	// Replace single quotes with '\''
	escaped := strings.ReplaceAll(arg, "'", "'\\''")
	
	// Wrap in single quotes
	return "'" + escaped + "'"
}
