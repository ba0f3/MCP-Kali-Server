package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/ba0f3/MCP-Kali-Server/pkg/handlers"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)


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
	}, handlers.NmapScanHandler)

	// Add Gobuster tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "gobuster_scan",
		Description: "Execute Gobuster to find directories, DNS subdomains, or virtual hosts",
	}, handlers.GobusterScanHandler)

	// Add Dirb tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "dirb_scan",
		Description: "Execute Dirb web content scanner",
	}, handlers.DirbScanHandler)

	// Add Nikto tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "nikto_scan",
		Description: "Execute Nikto web server scanner",
	}, handlers.NiktoScanHandler)

	// Add SQLmap tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "sqlmap_scan",
		Description: "Execute SQLmap SQL injection scanner",
	}, handlers.SqlmapScanHandler)

	// Add Hydra tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "hydra_attack",
		Description: "Execute Hydra password cracking tool",
	}, handlers.HydraAttackHandler)

	// Add John the Ripper tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "john_crack",
		Description: "Execute John the Ripper password cracker",
	}, handlers.JohnCrackHandler)

	// Add WPScan tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "wpscan_analyze",
		Description: "Execute WPScan WordPress vulnerability scanner",
	}, handlers.WpscanAnalyzeHandler)

	// Add Enum4linux tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "enum4linux_scan",
		Description: "Execute Enum4linux Windows/Samba enumeration tool",
	}, handlers.Enum4linuxScanHandler)

	// Add generic command execution tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "execute_command",
		Description: "Execute an arbitrary command on the Kali server",
	}, handlers.ExecuteCommandHandler)

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
