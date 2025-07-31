package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ba0f3/MCP-Kali-Server/pkg/executor"
	"github.com/ba0f3/MCP-Kali-Server/pkg/handlers"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)


func main() {
	var (
		debug = flag.Bool("debug", false, "Enable debug logging")
		httpAddr = flag.String("http", "", "if set, use streamable HTTP at this address, instead of stdin/stdout")
		timeout = flag.Int("timeout", 180, "Command execution timeout in seconds")
	)
	flag.Parse()

	// Set the global command timeout
	executor.SetGlobalTimeout(time.Duration(*timeout) * time.Second)

	// Print server configuration
	log.Println("=== MCP-Kali-Server Configuration ===")
	log.Println("Server Mode: MCP")
	log.Printf("Command Timeout: %d seconds", *timeout)
	log.Printf("Debug Mode: %v", *debug)
	if *httpAddr != "" {
		log.Printf("HTTP Address: %s", *httpAddr)
	} else {
		log.Println("Transport: stdio")
	}
	log.Println("Available Tools:")
	log.Println("  - nmap_scan")
	log.Println("  - gobuster_scan")
	log.Println("  - dirb_scan")
	log.Println("  - nikto_scan")
	log.Println("  - sqlmap_scan")
	log.Println("  - hydra_attack")
	log.Println("  - john_crack")
	log.Println("  - wpscan_analyze")
	log.Println("  - enum4linux_scan")
	log.Println("  - execute_command")
	log.Println("=====================================")

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



	if *httpAddr != "" {
		handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
			return server
		}, nil)
		log.Printf("Starting MCP Server with Kali Linux tools and listening at %s", *httpAddr)
		if err := http.ListenAndServe(*httpAddr, handler); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	} else {
		log.Println("Starting MCP Server with Kali Linux tools...")
		// Use stdio transport for MCP communication
		t := mcp.NewLoggingTransport(mcp.NewStdioTransport(), os.Stderr)
		if err := server.Run(context.Background(), t); err != nil {
			log.Printf("Server failed: %v", err)
		}
	}
}
