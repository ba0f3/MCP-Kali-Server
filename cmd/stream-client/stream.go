package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// StreamingMCPClient extends MCPClient with streaming capabilities
type StreamingMCPClient struct {
	*MCPClient
	streamTimeout time.Duration
}

// NewStreamingMCPClient creates a new streaming MCP client
func NewStreamingMCPClient(serverURL string, debug bool) *StreamingMCPClient {
	return &StreamingMCPClient{
		MCPClient:     NewMCPClient(serverURL, debug),
		streamTimeout: 10 * time.Minute,
	}
}

// CallToolStream calls a tool and streams the response
func (c *StreamingMCPClient) CallToolStream(name string, arguments interface{}, callback func(content ContentItem)) error {
	c.requestID++
	
	var argsJSON json.RawMessage
	if arguments != nil {
		data, err := json.Marshal(arguments)
		if err != nil {
			return fmt.Errorf("failed to marshal arguments: %v", err)
		}
		argsJSON = data
	}
	
	params := CallToolParams{
		Name:      name,
		Arguments: argsJSON,
	}
	
	var paramsJSON json.RawMessage
	data, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal params: %v", err)
	}
	paramsJSON = data
	
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params:  paramsJSON,
		ID:      c.requestID,
	}
	
	if c.debug {
		reqJSON, _ := json.MarshalIndent(request, "", "  ")
		fmt.Printf("Sending streaming request:\n%s\n", string(reqJSON))
	}
	
	reqBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}
	
	// Create request with streaming headers
	req, err := http.NewRequest("POST", c.serverURL, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	
	// Use a client without timeout for streaming
	streamClient := &http.Client{
		Timeout: 0, // No timeout for streaming
	}
	
	resp, err := streamClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}
	
	// Check if response is streaming
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/event-stream") {
		return c.handleSSEStream(resp.Body, callback)
	} else {
		// Fall back to regular response handling
		var response JSONRPCResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return fmt.Errorf("failed to decode response: %v", err)
		}
		
		if response.Error != nil {
			return fmt.Errorf("tool call error: %s", response.Error.Message)
		}
		
		var result CallToolResult
		if err := json.Unmarshal(response.Result, &result); err != nil {
			return fmt.Errorf("failed to unmarshal tool result: %v", err)
		}
		
		// Call callback for each content item
		for _, content := range result.Content {
			callback(content)
		}
		
		return nil
	}
}

// handleSSEStream processes Server-Sent Events stream
func (c *StreamingMCPClient) handleSSEStream(body io.Reader, callback func(ContentItem)) error {
	scanner := bufio.NewScanner(body)
	var eventData strings.Builder
	
	for scanner.Scan() {
		line := scanner.Text()
		
		if line == "" {
			// Empty line indicates end of event
			if eventData.Len() > 0 {
				if err := c.processSSEEvent(eventData.String(), callback); err != nil {
					if c.debug {
						fmt.Printf("Error processing SSE event: %v\n", err)
					}
				}
				eventData.Reset()
			}
			continue
		}
		
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			eventData.WriteString(data)
		}
	}
	
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stream: %v", err)
	}
	
	return nil
}

// processSSEEvent processes a single SSE event
func (c *StreamingMCPClient) processSSEEvent(data string, callback func(ContentItem)) error {
	// Handle special SSE messages
	if data == "[DONE]" {
		return nil
	}
	
	// Try to parse as content item
	var content ContentItem
	if err := json.Unmarshal([]byte(data), &content); err != nil {
		// Try to parse as a response wrapper
		var wrapper struct {
			Type    string          `json:"type"`
			Content json.RawMessage `json:"content"`
		}
		
		if err := json.Unmarshal([]byte(data), &wrapper); err != nil {
			return fmt.Errorf("failed to parse SSE data: %v", err)
		}
		
		// Extract content based on type
		switch wrapper.Type {
		case "content", "text":
			if err := json.Unmarshal(wrapper.Content, &content); err != nil {
				// Fallback to text content
				content = ContentItem{
					Type: "text",
					Text: string(wrapper.Content),
				}
			}
		default:
			// Unknown type, skip
			return nil
		}
	}
	
	callback(content)
	return nil
}

// StreamingDemo demonstrates streaming capabilities
func StreamingDemo(serverURL string) {
	fmt.Println("=== Streaming MCP Client Demo ===")
	fmt.Println()
	
	client := NewStreamingMCPClient(serverURL, false)
	
	// Initialize
	fmt.Println("Initializing connection...")
	initResult, err := client.Initialize()
	if err != nil {
		fmt.Printf("Failed to initialize: %v\n", err)
		return
	}
	
	fmt.Printf("Connected to: %s\n", initResult.ServerInfo.Name)
	fmt.Println()
	
	// Example 1: Stream a long-running command
	fmt.Println("Example 1: Streaming output from a long-running command")
	fmt.Println("Running: ping -n 5 google.com")
	fmt.Println("Output:")
	
	err = client.CallToolStream("execute_command", map[string]string{
		"command": "ping -n 5 google.com",
	}, func(content ContentItem) {
		if content.Type == "text" {
			// Print each line as it arrives
			lines := strings.Split(strings.TrimSpace(content.Text), "\n")
			for _, line := range lines {
				if line != "" {
					fmt.Printf("  > %s\n", line)
				}
			}
		}
	})
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	
	fmt.Println()
	
	// Example 2: Stream nmap output
	fmt.Println("Example 2: Streaming nmap scan output")
	fmt.Println("Scanning localhost...")
	
	startTime := time.Now()
	lineCount := 0
	
	err = client.CallToolStream("nmap_scan", map[string]interface{}{
		"target":    "127.0.0.1",
		"scan_type": "-sV",
		"ports":     "1-1000",
	}, func(content ContentItem) {
		if content.Type == "text" {
			lines := strings.Split(strings.TrimSpace(content.Text), "\n")
			for _, line := range lines {
				if line != "" {
					lineCount++
					elapsed := time.Since(startTime).Seconds()
					fmt.Printf("[%.1fs] %s\n", elapsed, line)
				}
			}
		}
	})
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	
	fmt.Printf("\nTotal lines received: %d\n", lineCount)
	fmt.Printf("Total time: %.1fs\n", time.Since(startTime).Seconds())
}

// MonitorMode runs continuous monitoring with streaming
func MonitorMode(serverURL string, command string, interval time.Duration) {
	fmt.Printf("=== Monitor Mode ===\n")
	fmt.Printf("Command: %s\n", command)
	fmt.Printf("Interval: %v\n", interval)
	fmt.Printf("Press Ctrl+C to stop\n\n")
	
	client := NewStreamingMCPClient(serverURL, false)
	
	// Initialize
	if _, err := client.Initialize(); err != nil {
		fmt.Printf("Failed to initialize: %v\n", err)
		return
	}
	
	iteration := 0
	for {
		iteration++
		fmt.Printf("\n[Iteration %d - %s]\n", iteration, time.Now().Format("15:04:05"))
		
		err := client.CallToolStream("execute_command", map[string]string{
			"command": command,
		}, func(content ContentItem) {
			if content.Type == "text" {
				fmt.Print(content.Text)
			}
		})
		
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		
		time.Sleep(interval)
	}
}
