package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// JSON-RPC 2.0 structures
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      interface{}     `json:"id"`
}

type JSONRPCResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	Result  json.RawMessage  `json:"result,omitempty"`
	Error   *JSONRPCError    `json:"error,omitempty"`
	ID      interface{}      `json:"id"`
}

type JSONRPCError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// MCP specific structures
type InitializeParams struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    ClientCapabilities     `json:"capabilities"`
	ClientInfo      ClientInfo             `json:"clientInfo"`
}

type ClientCapabilities struct {
	Sampling    json.RawMessage `json:"sampling,omitempty"`
	Experimental json.RawMessage `json:"experimental,omitempty"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResult struct {
	ProtocolVersion string           `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      ServerInfo        `json:"serverInfo"`
}

type ServerCapabilities struct {
	Tools        *ToolsCapability        `json:"tools,omitempty"`
	Resources    json.RawMessage         `json:"resources,omitempty"`
	Prompts      json.RawMessage         `json:"prompts,omitempty"`
	Logging      json.RawMessage         `json:"logging,omitempty"`
	Experimental json.RawMessage         `json:"experimental,omitempty"`
}

type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

type ListToolsResult struct {
	Tools []Tool `json:"tools"`
}

type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

type CallToolParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments,omitempty"`
}

type CallToolResult struct {
	Content []ContentItem `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

type ContentItem struct {
	Type     string          `json:"type"`
	Text     string          `json:"text,omitempty"`
	Data     json.RawMessage `json:"data,omitempty"`
	MimeType string          `json:"mimeType,omitempty"`
}

// MCPClient represents a streamable HTTP MCP client
type MCPClient struct {
	serverURL  string
	httpClient *http.Client
	requestID  int
	debug      bool
	sessionID  string
	authType   string
	authSecret string
}

// NewMCPClient creates a new MCP client
func NewMCPClient(serverURL string, debug bool) *MCPClient {
	// Get auth configuration from environment
	authType := os.Getenv("AUTH_TYPE")
	authSecret := os.Getenv("AUTH_SECRET")
	
	// Default to apikey if authSecret is set but authType is not
	if authSecret != "" && authType == "" {
		authType = "apikey"
	}
	
	client := &MCPClient{
		serverURL: serverURL,
		httpClient: &http.Client{
			Timeout: 300 * time.Second, // 5 minute timeout for long-running operations
		},
		requestID:  0,
		debug:      debug,
		authType:   authType,
		authSecret: authSecret,
	}
	
	if debug && authSecret != "" {
		log.Printf("Authentication configured: %s", authType)
	}
	
	return client
}

// sendNotification sends a JSON-RPC notification (no response expected)
func (c *MCPClient) sendNotification(method string, params interface{}) error {
	var paramsJSON json.RawMessage
	if params != nil {
		data, err := json.Marshal(params)
		if err != nil {
			return fmt.Errorf("failed to marshal params: %v", err)
		}
		paramsJSON = data
	}
	
	notification := struct {
		JSONRPC string          `json:"jsonrpc"`
		Method  string          `json:"method"`
		Params  json.RawMessage `json:"params,omitempty"`
	}{
		JSONRPC: "2.0",
		Method:  method,
		Params:  paramsJSON,
	}
	
	if c.debug {
		reqJSON, _ := json.MarshalIndent(notification, "", "  ")
		log.Printf("Sending notification:\n%s\n", string(reqJSON))
	}
	
	reqBody, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %v", err)
	}
	
	req, err := http.NewRequest("POST", c.serverURL, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")

	// Set authentication header if configured
	if c.authSecret != "" {
		switch c.authType {
		case "apikey":
			req.Header.Set("X-API-Key", c.authSecret)
		case "bearer":
			req.Header.Set("Authorization", "Bearer " + c.authSecret)
		}
	}
	
	// Add session ID if available
	if c.sessionID != "" {
		req.Header.Set("Mcp-Session-Id", c.sessionID)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send notification: %v", err)
	}
	defer resp.Body.Close()
	
	if c.debug {
		log.Printf("Notification response headers: %v", resp.Header)
	}
	
	// Capture session ID from response headers
	if sessionID := resp.Header.Get("Mcp-Session-Id"); sessionID != "" && c.sessionID == "" {
		c.sessionID = sessionID
		if c.debug {
			log.Printf("Captured session ID: %s", c.sessionID)
		}
	}
	
	// For notifications, we don't expect a response body
	// Just check the status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}
	
	return nil
}

// sendRequest sends a JSON-RPC request and returns the response
func (c *MCPClient) sendRequest(method string, params interface{}) (*JSONRPCResponse, error) {
	c.requestID++
	
	var paramsJSON json.RawMessage
	if params != nil {
		data, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal params: %v", err)
		}
		paramsJSON = data
	}
	
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  paramsJSON,
		ID:      c.requestID,
	}
	
	if c.debug {
		reqJSON, _ := json.MarshalIndent(request, "", "  ")
		log.Printf("Sending request:\n%s\n", string(reqJSON))
	}
	
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}
	
	req, err := http.NewRequest("POST", c.serverURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
req.Header.Set("Accept", "application/json, text/event-stream")

	// Set authentication header if configured
	if c.authSecret != "" {
		switch c.authType {
		case "apikey":
			req.Header.Set("X-API-Key", c.authSecret)
		case "bearer":
			req.Header.Set("Authorization", "Bearer " + c.authSecret)
		}
	}
	// Add session ID if available
	if c.sessionID != "" {
		req.Header.Set("Mcp-Session-Id", c.sessionID)
		if c.debug {
			log.Printf("Using session ID: %s", c.sessionID)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	
	// Capture session ID from response headers
	if sessionID := resp.Header.Get("Mcp-Session-Id"); sessionID != "" && c.sessionID == "" {
		c.sessionID = sessionID
		if c.debug {
			log.Printf("Captured session ID from request: %s", c.sessionID)
		}
	}
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}
	
	// Check if response is SSE formatted
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/event-stream") {
		// Handle SSE formatted response
		scanner := bufio.NewScanner(resp.Body)
		var jsonData string
		
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "data: ") {
				jsonData = strings.TrimPrefix(line, "data: ")
				break // We only need the first data line for non-streaming responses
			}
		}
		
		if jsonData == "" {
			return nil, fmt.Errorf("no data found in SSE response")
		}
		
		var response JSONRPCResponse
		if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
			return nil, fmt.Errorf("failed to decode SSE response: %v", err)
		}
		return &response, nil
	} else {
		// Handle regular JSON response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}
		
		// If the body looks like SSE even without the header, parse it
		if strings.HasPrefix(string(body), "event:") {
			lines := strings.Split(string(body), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "data: ") {
					jsonData := strings.TrimPrefix(line, "data: ")
					var response JSONRPCResponse
					if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
						return nil, fmt.Errorf("failed to decode SSE response: %v", err)
					}
					return &response, nil
				}
			}
			return nil, fmt.Errorf("no data found in SSE response")
		}
		
		var response JSONRPCResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("failed to decode response: %v (body: %s)", err, string(body))
		}
		return &response, nil
	}
}

// Initialize performs the MCP initialization handshake
func (c *MCPClient) Initialize() (*InitializeResult, error) {
	params := InitializeParams{
		ProtocolVersion: "2024-11-05",
		Capabilities:    ClientCapabilities{},
		ClientInfo: ClientInfo{
			Name:    "mcp-stream-client",
			Version: "1.0.0",
		},
	}
	
	resp, err := c.sendRequest("initialize", params)
	if err != nil {
		return nil, err
	}
	
	if resp.Error != nil {
		return nil, fmt.Errorf("initialization error: %s", resp.Error.Message)
	}
	
	var result InitializeResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal initialize result: %v", err)
	}
	
	// Send initialized notification
	err = c.sendNotification("initialized", struct{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to send initialized notification: %v", err)
	}
	
	return &result, nil
}

// ListTools retrieves the list of available tools
func (c *MCPClient) ListTools() ([]Tool, error) {
	resp, err := c.sendRequest("tools/list", nil)
	if err != nil {
		return nil, err
	}
	
	if resp.Error != nil {
		return nil, fmt.Errorf("list tools error: %s", resp.Error.Message)
	}
	
	var result ListToolsResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tools list: %v", err)
	}
	
	return result.Tools, nil
}

// CallTool calls a specific tool with arguments
func (c *MCPClient) CallTool(name string, arguments interface{}) (*CallToolResult, error) {
	var argsJSON json.RawMessage
	if arguments != nil {
		data, err := json.Marshal(arguments)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal arguments: %v", err)
		}
		argsJSON = data
	}
	
	params := CallToolParams{
		Name:      name,
		Arguments: argsJSON,
	}
	
	resp, err := c.sendRequest("tools/call", params)
	if err != nil {
		return nil, err
	}
	
	if resp.Error != nil {
		return nil, fmt.Errorf("tool call error: %s", resp.Error.Message)
	}
	
	var result CallToolResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tool result: %v", err)
	}
	
	return &result, nil
}

// Interactive mode
func (c *MCPClient) runInteractive() error {
	fmt.Println("MCP Stream Client - Interactive Mode")
	fmt.Println("Commands:")
	fmt.Println("  list              - List available tools")
	fmt.Println("  call <tool> <args> - Call a tool with JSON arguments")
	fmt.Println("  exit              - Exit the client")
	fmt.Println()
	
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		
		parts := strings.Fields(line)
		command := parts[0]
		
		switch command {
		case "exit", "quit":
			return nil
			
		case "list":
			tools, err := c.ListTools()
			if err != nil {
				fmt.Printf("Error listing tools: %v\n", err)
				continue
			}
			fmt.Println("Available tools:")
			for _, tool := range tools {
				fmt.Printf("  - %s: %s\n", tool.Name, tool.Description)
			}
			
		case "call":
			if len(parts) < 2 {
				fmt.Println("Usage: call <tool> [<json-args>]")
				continue
			}
			
			toolName := parts[1]
			var args interface{}
			
			if len(parts) > 2 {
				argsStr := strings.Join(parts[2:], " ")
				if err := json.Unmarshal([]byte(argsStr), &args); err != nil {
					fmt.Printf("Error parsing arguments: %v\n", err)
					continue
				}
			}
			
			result, err := c.CallTool(toolName, args)
			if err != nil {
				fmt.Printf("Error calling tool: %v\n", err)
				continue
			}
			
			if result.IsError {
				fmt.Println("Tool returned an error:")
			} else {
				fmt.Println("Tool result:")
			}
			
			for _, content := range result.Content {
				if content.Type == "text" {
					fmt.Println(content.Text)
				} else {
					fmt.Printf("Content type: %s\n", content.Type)
					if content.Data != nil {
						fmt.Printf("Data: %s\n", string(content.Data))
					}
				}
			}
			
		default:
			fmt.Printf("Unknown command: %s\n", command)
		}
	}
	
	return scanner.Err()
}

// Test mode - run predefined tests
func (c *MCPClient) runTests() error {
	fmt.Println("Running MCP client tests...")
	fmt.Println()
	
	// Test 1: List tools
	fmt.Println("Test 1: Listing tools")
	tools, err := c.ListTools()
	if err != nil {
		return fmt.Errorf("failed to list tools: %v", err)
	}
	
	fmt.Printf("Found %d tools:\n", len(tools))
	for _, tool := range tools {
		fmt.Printf("  - %s: %s\n", tool.Name, tool.Description)
	}
	fmt.Println()
	
	// Test 2: Execute a simple command
	fmt.Println("Test 2: Execute simple command")
	result, err := c.CallTool("execute_command", map[string]string{
		"command": "whoami",
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Result:")
		for _, content := range result.Content {
			if content.Type == "text" {
				fmt.Println(content.Text)
			}
		}
	}
	fmt.Println()
	
	// Test 3: Nmap scan (localhost)
	fmt.Println("Test 3: Nmap scan on localhost")
	result, err = c.CallTool("nmap_scan", map[string]interface{}{
		"target":    "127.0.0.1",
		"scan_type": "-sV",
		"ports":     "22,80,443",
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Result:")
		for _, content := range result.Content {
			if content.Type == "text" {
				// Limit output for readability
				lines := strings.Split(content.Text, "\n")
				for i, line := range lines {
					if i > 20 {
						fmt.Println("... (truncated)")
						break
					}
					fmt.Println(line)
				}
			}
		}
	}
	
	return nil
}

func main() {
	var (
		serverURL   = flag.String("server", "http://localhost:8080", "MCP server URL")
		interactive = flag.Bool("interactive", false, "Run in interactive mode")
		test        = flag.Bool("test", false, "Run test suite")
		debug       = flag.Bool("debug", false, "Enable debug logging")
		toolName    = flag.String("tool", "", "Tool to call")
		toolArgs    = flag.String("args", "", "Tool arguments as JSON")
		stream      = flag.Bool("stream", false, "Run streaming demo")
		monitor     = flag.String("monitor", "", "Monitor mode: run command repeatedly")
		interval    = flag.Duration("interval", 5*time.Second, "Monitor interval")
	)
	flag.Parse()
	
	client := NewMCPClient(*serverURL, *debug)
	
	// Initialize the connection
	fmt.Printf("Connecting to MCP server at %s...\n", *serverURL)
	initResult, err := client.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}
	
	fmt.Printf("Connected to: %s\n", initResult.ServerInfo.Name)
	if initResult.ServerInfo.Version != "" {
		fmt.Printf("Server version: %s\n", initResult.ServerInfo.Version)
	}
	fmt.Printf("Protocol version: %s\n", initResult.ProtocolVersion)
	fmt.Println()
	
	// Run the appropriate mode
	switch {
	case *interactive:
		if err := client.runInteractive(); err != nil {
			log.Fatalf("Interactive mode error: %v", err)
		}
		
	case *test:
		if err := client.runTests(); err != nil {
			log.Fatalf("Test mode error: %v", err)
		}
		
	case *stream:
		StreamingDemo(*serverURL)
		
	case *monitor != "":
		MonitorMode(*serverURL, *monitor, *interval)
		
	case *toolName != "":
		// Single tool call mode
		var args interface{}
		if *toolArgs != "" {
			if err := json.Unmarshal([]byte(*toolArgs), &args); err != nil {
				log.Fatalf("Failed to parse arguments: %v", err)
			}
		}
		
		result, err := client.CallTool(*toolName, args)
		if err != nil {
			log.Fatalf("Failed to call tool: %v", err)
		}
		
		for _, content := range result.Content {
			if content.Type == "text" {
				fmt.Println(content.Text)
			} else {
				fmt.Printf("Content type: %s\n", content.Type)
				if content.Data != nil {
					fmt.Printf("Data: %s\n", string(content.Data))
				}
			}
		}
		
	default:
		// Default: list tools
		tools, err := client.ListTools()
		if err != nil {
			log.Fatalf("Failed to list tools: %v", err)
		}
		
		fmt.Println("Available tools:")
		for _, tool := range tools {
			fmt.Printf("  - %s: %s\n", tool.Name, tool.Description)
		}
		fmt.Println("\nUse -interactive for interactive mode or -test to run tests")
	}
}
