package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
)

// StreamEvent represents a server-sent event
type StreamEvent struct {
	Type      string    `json:"type"`      // "stdout", "stderr", "exit"
	Data      string    `json:"data"`      // The actual output line
	Timestamp time.Time `json:"timestamp"` // When the event occurred
	ExitCode  int       `json:"exit_code,omitempty"` // Only for "exit" type
}

// StreamCommandHandler executes a command and streams the output
func StreamCommandHandler(c *gin.Context) {
	var data map[string]string
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
		return
	}

	command, ok := data["command"]
	if !ok || command == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Command parameter is required"})
		return
	}

	// Set headers for streaming
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	// Create command with context for cancellation
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", command)

	// Get stdout and stderr pipes
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		sendSSEError(c, fmt.Sprintf("Failed to create stdout pipe: %v", err))
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		sendSSEError(c, fmt.Sprintf("Failed to create stderr pipe: %v", err))
		return
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		sendSSEError(c, fmt.Sprintf("Failed to start command: %v", err))
		return
	}

	// Channel to signal when readers are done
	done := make(chan bool, 2)

	// Read stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			event := StreamEvent{
				Type:      "stdout",
				Data:      scanner.Text(),
				Timestamp: time.Now(),
			}
			sendSSEEvent(c, event)
		}
		done <- true
	}()

	// Read stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			event := StreamEvent{
				Type:      "stderr",
				Data:      scanner.Text(),
				Timestamp: time.Now(),
			}
			sendSSEEvent(c, event)
		}
		done <- true
	}()

	// Wait for both readers to finish
	<-done
	<-done

	// Wait for command to complete
	err = cmd.Wait()
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		}
	}

	// Send exit event
	event := StreamEvent{
		Type:      "exit",
		Data:      "Command completed",
		Timestamp: time.Now(),
		ExitCode:  exitCode,
	}
	sendSSEEvent(c, event)

	// Flush any remaining data
	c.Writer.Flush()
}

// sendSSEEvent sends a server-sent event
func sendSSEEvent(c *gin.Context, event StreamEvent) {
	data, _ := json.Marshal(event)
	fmt.Fprintf(c.Writer, "data: %s\n\n", string(data))
	c.Writer.Flush()
}

// sendSSEError sends an error as a server-sent event
func sendSSEError(c *gin.Context, errMsg string) {
	event := StreamEvent{
		Type:      "error",
		Data:      errMsg,
		Timestamp: time.Now(),
	}
	sendSSEEvent(c, event)
}
