package executor

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	// DefaultTimeout is the default command execution timeout
	DefaultTimeout = 15 * time.Minute
)

var (
	// GlobalTimeout is the global command execution timeout that can be configured
	GlobalTimeout = DefaultTimeout
)

// Result represents the result of a command execution
type Result struct {
	Stdout       string `json:"stdout"`
	Stderr       string `json:"stderr"`
	ReturnCode   int    `json:"return_code"`
	Success      bool   `json:"success"`
	TimedOut     bool   `json:"timed_out"`
	PartialResults bool `json:"partial_results"`
}

// CommandExecutor handles command execution with timeout management
type CommandExecutor struct {
	Command   string
	Timeout   time.Duration
	stdout    strings.Builder
	stderr    strings.Builder
	returnCode int
	timedOut  bool
	mu        sync.Mutex
}

// NewCommandExecutor creates a new CommandExecutor
func NewCommandExecutor(command string, timeout time.Duration) *CommandExecutor {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &CommandExecutor{
		Command: command,
		Timeout: timeout,
	}
}

// Execute runs the command and returns the result
func (ce *CommandExecutor) Execute() (*Result, error) {
	log.Printf("Executing command: %s", ce.Command)

	ctx, cancel := context.WithTimeout(context.Background(), ce.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", ce.Command)

	// Create pipes for stdout and stderr
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	// Create wait group for goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	// Read stdout
	go func() {
		defer wg.Done()
		ce.readPipe(stdoutPipe, &ce.stdout)
	}()

	// Read stderr
	go func() {
		defer wg.Done()
		ce.readPipe(stderrPipe, &ce.stderr)
	}()

	// Wait for command to complete
	cmdErr := cmd.Wait()

	// Wait for all output to be read
	wg.Wait()

	// Check if command timed out
	if ctx.Err() == context.DeadlineExceeded {
		ce.timedOut = true
		ce.returnCode = -1
		log.Printf("Command timed out after %v", ce.Timeout)
	} else if cmdErr != nil {
		if exitError, ok := cmdErr.(*exec.ExitError); ok {
			ce.returnCode = exitError.ExitCode()
		} else {
			ce.returnCode = -1
		}
	} else {
		ce.returnCode = 0
	}

	// Determine success
	success := ce.returnCode == 0
	if ce.timedOut && (ce.stdout.Len() > 0 || ce.stderr.Len() > 0) {
		success = true // Consider it a success if we have output even with timeout
	}

	return &Result{
		Stdout:         ce.stdout.String(),
		Stderr:         ce.stderr.String(),
		ReturnCode:     ce.returnCode,
		Success:        success,
		TimedOut:       ce.timedOut,
		PartialResults: ce.timedOut && (ce.stdout.Len() > 0 || ce.stderr.Len() > 0),
	}, nil
}

// readPipe reads from a pipe and writes to a string builder
func (ce *CommandExecutor) readPipe(pipe io.Reader, builder *strings.Builder) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		ce.mu.Lock()
		builder.WriteString(scanner.Text())
		builder.WriteString("\n")
		ce.mu.Unlock()
	}
}

// ExecuteCommand is a convenience function to execute a command
func ExecuteCommand(command string) (*Result, error) {
	executor := NewCommandExecutor(command, GlobalTimeout)
	return executor.Execute()
}

// SetGlobalTimeout sets the global command execution timeout
func SetGlobalTimeout(timeout time.Duration) {
	GlobalTimeout = timeout
}
