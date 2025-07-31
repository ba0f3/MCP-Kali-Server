package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Wh0am123/MCP-Kali-Server/pkg/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Tool handler functions
func NmapScanHandler(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[tools.NmapParams]) (*mcp.CallToolResultFor[any], error) {
	result, err := tools.NmapScan(params.Arguments)
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: formatToolResult(result)},
		},
	}, nil
}

func GobusterScanHandler(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[tools.GobusterParams]) (*mcp.CallToolResultFor[any], error) {
	result, err := tools.GobusterScan(params.Arguments)
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: formatToolResult(result)},
		},
	}, nil
}

func DirbScanHandler(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[tools.DirbParams]) (*mcp.CallToolResultFor[any], error) {
	result, err := tools.DirbScan(params.Arguments)
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: formatToolResult(result)},
		},
	}, nil
}

func NiktoScanHandler(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[tools.NiktoParams]) (*mcp.CallToolResultFor[any], error) {
	result, err := tools.NiktoScan(params.Arguments)
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: formatToolResult(result)},
		},
	}, nil
}

func SqlmapScanHandler(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[tools.SqlmapParams]) (*mcp.CallToolResultFor[any], error) {
	result, err := tools.SqlmapScan(params.Arguments)
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: formatToolResult(result)},
		},
	}, nil
}

func HydraAttackHandler(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[tools.HydraParams]) (*mcp.CallToolResultFor[any], error) {
	result, err := tools.HydraAttack(params.Arguments)
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: formatToolResult(result)},
		},
	}, nil
}

func JohnCrackHandler(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[tools.JohnParams]) (*mcp.CallToolResultFor[any], error) {
	result, err := tools.JohnCrack(params.Arguments)
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: formatToolResult(result)},
		},
	}, nil
}

func WpscanAnalyzeHandler(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[tools.WpscanParams]) (*mcp.CallToolResultFor[any], error) {
	result, err := tools.WpscanAnalyze(params.Arguments)
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: formatToolResult(result)},
		},
	}, nil
}

func Enum4linuxScanHandler(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[tools.Enum4linuxParams]) (*mcp.CallToolResultFor[any], error) {
	result, err := tools.Enum4linuxScan(params.Arguments)
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: formatToolResult(result)},
		},
	}, nil
}

func ExecuteCommandHandler(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[tools.GenericCommandParams]) (*mcp.CallToolResultFor[any], error) {
	result, err := tools.ExecuteGenericCommand(params.Arguments)
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: formatToolResult(result)},
		},
	}, nil
}

// formatToolResult formats the tool result for display
func formatToolResult(result *tools.ToolResult) string {
	if result.Success {
		return result.Stdout
	}
	return fmt.Sprintf("Error: %s\n%s", result.Error, result.Stderr)
}

func main() {
	var (
		debug = flag.Bool("debug", false, "Enable debug logging")
	)
	flag.Parse()

	if *debug {
		log.Println("Debug mode enabled")
	}

	// Initialize MCP server
	server := mcp.NewServer(&mcp.Implementation{Name: "kali-tools"}, nil)

	// Add Nmap tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "nmap_scan",
		Description: "Execute an Nmap scan against a target",
	}, NmapScanHandler)

	// Add Gobuster tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "gobuster_scan",
		Description: "Execute Gobuster to find directories, DNS subdomains, or virtual hosts",
	}, GobusterScanHandler)

	// Add Dirb tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "dirb_scan",
		Description: "Execute Dirb web content scanner",
	}, DirbScanHandler)

	// Add Nikto tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "nikto_scan",
		Description: "Execute Nikto web server scanner",
	}, NiktoScanHandler)

	// Add SQLmap tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "sqlmap_scan",
		Description: "Execute SQLmap SQL injection scanner",
	}, SqlmapScanHandler)

	// Add Hydra tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "hydra_attack",
		Description: "Execute Hydra password cracking tool",
	}, HydraAttackHandler)

	// Add John the Ripper tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "john_crack",
		Description: "Execute John the Ripper password cracker",
	}, JohnCrackHandler)

	// Add WPScan tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "wpscan_analyze",
		Description: "Execute WPScan WordPress vulnerability scanner",
	}, WpscanAnalyzeHandler)

	// Add Enum4linux tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "enum4linux_scan",
		Description: "Execute Enum4linux Windows/Samba enumeration tool",
	}, Enum4linuxScanHandler)

	// Add generic command execution tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "execute_command",
		Description: "Execute an arbitrary command on the Kali server",
	}, ExecuteCommandHandler)

	log.Println("Starting MCP Server with Kali Linux tools...")

	// Use stdio transport for MCP communication
	var t mcp.Transport
	t = mcp.NewStdioTransport()
	if *debug {
		t = mcp.NewLoggingTransport(t, os.Stderr)
	}

	// Start the MCP server
	if err := server.Run(context.Background(), t); err != nil {
		log.Fatalf("Failed to start MCP server: %v", err)
	}
}
