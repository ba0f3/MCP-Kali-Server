package handlers

import (
	"context"
	"fmt"

	"github.com/ba0f3/MCP-Kali-Server/pkg/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// NmapScanHandler handles Nmap scan requests
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

// GobusterScanHandler handles Gobuster scan requests
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

// DirbScanHandler handles Dirb scan requests
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

// NiktoScanHandler handles Nikto scan requests
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

// SqlmapScanHandler handles SQLmap scan requests
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

// HydraAttackHandler handles Hydra attack requests
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

// JohnCrackHandler handles John the Ripper requests
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

// WpscanAnalyzeHandler handles WPScan requests
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

// Enum4linuxScanHandler handles Enum4linux scan requests
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

// PingHandler handles ping requests
func PingHandler(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[tools.PingParams]) (*mcp.CallToolResultFor[any], error) {
	result, err := tools.Ping(params.Arguments)
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: formatToolResult(result)},
		},
	}, nil
}

// NucleiScanHandler handles Nuclei scan requests
func NucleiScanHandler(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[tools.NucleiParams]) (*mcp.CallToolResultFor[any], error) {
	result, err := tools.NucleiScan(params.Arguments)
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: formatToolResult(result)},
		},
	}, nil
}

// ExecuteCommandHandler handles generic command execution requests
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

// InitializeServer initializes the MCP server with tools
func InitializeServer() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "kali-tools"}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "nmap_scan",
		Description: "Execute an Nmap scan against a target",
	}, NmapScanHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "gobuster_scan",
		Description: "Execute Gobuster to find directories, DNS subdomains, or virtual hosts",
	}, GobusterScanHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "dirb_scan",
		Description: "Execute Dirb web content scanner",
	}, DirbScanHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "nikto_scan",
		Description: "Execute Nikto web server scanner",
	}, NiktoScanHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "sqlmap_scan",
		Description: "Execute SQLmap SQL injection scanner",
	}, SqlmapScanHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "hydra_attack",
		Description: "Execute Hydra password cracking tool",
	}, HydraAttackHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "john_crack",
		Description: "Execute John the Ripper password cracker",
	}, JohnCrackHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "wpscan_analyze",
		Description: "Execute WPScan WordPress vulnerability scanner",
	}, WpscanAnalyzeHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "enum4linux_scan",
		Description: "Execute Enum4linux Windows/Samba enumeration tool",
	}, Enum4linuxScanHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "ping",
		Description: "Execute ping to test network connectivity",
	}, PingHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "nuclei_scan",
		Description: "Execute Nuclei template-based vulnerability scanner",
	}, NucleiScanHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "execute_command",
		Description: "Execute an arbitrary command on the Kali server",
	}, ExecuteCommandHandler)

	return server
}

// formatToolResult formats the tool result for display
func formatToolResult(result *tools.ToolResult) string {
	if result.Success {
		return result.Stdout
	}
	return fmt.Sprintf("Error: %s\n%s", result.Error, result.Stderr)
}
