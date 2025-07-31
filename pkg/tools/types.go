package tools

// ToolResult represents the result of a tool execution
type ToolResult struct {
	Stdout         string `json:"stdout"`
	Stderr         string `json:"stderr"`
	Success        bool   `json:"success"`
	Error          string `json:"error,omitempty"`
	ReturnCode     int    `json:"return_code"`
	TimedOut       bool   `json:"timed_out"`
	PartialResults bool   `json:"partial_results"`
}
