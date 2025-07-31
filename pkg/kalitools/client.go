package kalitools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	// DefaultRequestTimeout is the default timeout for API requests
	DefaultRequestTimeout = 5 * time.Minute
)

// Client represents a client for communicating with the Kali Linux Tools API Server
type Client struct {
	serverURL string
	client    *http.Client
}

// NewClient creates a new Kali Tools Client
func NewClient(serverURL string, timeout time.Duration) *Client {
	if timeout == 0 {
		timeout = DefaultRequestTimeout
	}

	return &Client{
		serverURL: serverURL,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Response represents a generic API response
type Response struct {
	Stdout  string `json:"stdout,omitempty"`
	Stderr  string `json:"stderr,omitempty"`
	Error   string `json:"error,omitempty"`
	Success bool   `json:"success"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

// SafeGet performs a GET request with optional query parameters
func (c *Client) SafeGet(endpoint string, params map[string]string) (*Response, error) {
	url := fmt.Sprintf("%s/%s", c.serverURL, endpoint)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters
	if params != nil {
		q := req.URL.Query()
		for key, value := range params {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	log.Printf("GET %s", req.URL.String())

	resp, err := c.client.Do(req)
	if err != nil {
		return &Response{
			Error:   fmt.Sprintf("Request failed: %v", err),
			Success: false,
		}, nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Response{
			Error:   fmt.Sprintf("Failed to read response: %v", err),
			Success: false,
		}, nil
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return &Response{
			Error:   fmt.Sprintf("Failed to parse response: %v", err),
			Success: false,
		}, nil
	}

	return &result, nil
}

// SafePost performs a POST request with JSON data
func (c *Client) SafePost(endpoint string, data interface{}) (*Response, error) {
	url := fmt.Sprintf("%s/%s", c.serverURL, endpoint)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	log.Printf("POST %s with data: %s", url, string(jsonData))

	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return &Response{
			Error:   fmt.Sprintf("Request failed: %v", err),
			Success: false,
		}, nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Response{
			Error:   fmt.Sprintf("Failed to read response: %v", err),
			Success: false,
		}, nil
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return &Response{
			Error:   fmt.Sprintf("Failed to parse response: %v", err),
			Success: false,
		}, nil
	}

	return &result, nil
}

// ExecuteCommand executes a generic command on the Kali server
func (c *Client) ExecuteCommand(command string) (*Response, error) {
	data := map[string]string{"command": command}
	return c.SafePost("api/command", data)
}

// CheckHealth checks the health of the Kali Tools API Server
func (c *Client) CheckHealth() (*Response, error) {
	return c.SafeGet("health", nil)
}
