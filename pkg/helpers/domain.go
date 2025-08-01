package helpers

import (
	"fmt"
	"net/url"
	"strings"
)

// ParseDomain extracts the domain (host) from a given input string.
// It can handle full URLs, strings without a scheme (e.g., "example.com/path"),
// or just domain names (e.g., "example.com").
func ParseDomain(rawURL string) (string, error) {
	// 1. Trim leading and trailing whitespace from the input string.
	trimmedURL := strings.TrimSpace(rawURL)
	if trimmedURL == "" {
		return "", fmt.Errorf("error: input URL string is empty")
	}

	// 2. If the string does not contain "://", prepend "http://".
	// This helps url.Parse to correctly handle cases like "example.com/path".
	parseableURL := trimmedURL
	if !strings.Contains(parseableURL, "://") {
		parseableURL = "http://" + parseableURL
	}

	// 3. Use the standard library to parse the URL.
	u, err := url.Parse(parseableURL)
	if err != nil {
		// Return an error if the string is invalid after prepending the scheme.
		return "", fmt.Errorf("error parsing URL '%s': %w", rawURL, err)
	}

	// 4. Check if the Host (domain) is empty.
	// This can happen for inputs like "/just/a/path".
	if u.Host == "" {
		return "", fmt.Errorf("error: could not extract domain from '%s'", rawURL)
	}

	// 5. Return the Host, which is the domain name.
	return u.Host, nil
}
