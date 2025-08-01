package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ba0f3/MCP-Kali-Server/pkg/executor"
	"github.com/ba0f3/MCP-Kali-Server/pkg/handlers"
	"github.com/ba0f3/MCP-Kali-Server/pkg/middleware"
	"github.com/ba0f3/MCP-Kali-Server/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// handleServiceInstall installs the server as a system service
func handleServiceInstall(serviceName, servicePort string) error {
	config := service.GetDefaultConfig()
	config.Name = serviceName

	// Set service arguments to use HTTP mode with the specified port
	config.Args = []string{
		"-http", servicePort,
		"-timeout", "900",
	}

	fmt.Printf("Installing %s service...\n", serviceName)
	fmt.Printf("Service will listen on port: %s\n", servicePort)
	fmt.Printf("Executable: %s\n", config.Executable)
	fmt.Printf("Working directory: %s\n", config.WorkingDir)

	if err := service.InstallService(config); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	fmt.Printf("\nService '%s' installed successfully!\n", serviceName)
	fmt.Println("\nTo start the service:")
	fmt.Printf("  Windows: sc start %s\n", serviceName)
	fmt.Printf("  Linux: sudo systemctl start %s\n", serviceName)
	fmt.Printf("  macOS: sudo launchctl start %s\n", serviceName)

	return nil
}

// handleServiceUninstall removes the installed service
func handleServiceUninstall(serviceName string) error {
	fmt.Printf("Uninstalling %s service...\n", serviceName)

	if err := service.UninstallService(serviceName); err != nil {
		return fmt.Errorf("uninstallation failed: %w", err)
	}

	fmt.Printf("Service '%s' uninstalled successfully!\n", serviceName)
	return nil
}

func main() {
	var (
		debug = flag.Bool("debug", false, "Enable debug logging")
		httpAddr = flag.String("http", "", "if set, use streamable HTTP at this address, instead of stdin/stdout")
		timeout = flag.Int("timeout", 900, "Command execution timeout in seconds")
		installService = flag.Bool("install-service", false, "Install the server as a system service")
		uninstallService = flag.Bool("uninstall-service", false, "Uninstall the system service")
		serviceName = flag.String("service-name", "mcp-kali-server", "Name of the service")
		servicePort = flag.String("service-port", ":8080", "Port for the service to listen on (used with -install-service)")
	)
	flag.Parse()

	// Handle service installation/uninstallation
	if *installService {
		if err := handleServiceInstall(*serviceName, *servicePort); err != nil {
			log.Fatalf("Failed to install service: %v", err)
		}
		return
	}

	if *uninstallService {
		if err := handleServiceUninstall(*serviceName); err != nil {
			log.Fatalf("Failed to uninstall service: %v", err)
		}
		return
	}

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
	log.Println("  - ping")
	log.Println("  - nuclei_scan")
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

	// Add Ping tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "ping",
		Description: "Execute ping to test network connectivity",
	}, handlers.PingHandler)

	// Add Nuclei tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "nuclei_scan",
		Description: "Execute Nuclei template-based vulnerability scanner",
	}, handlers.NucleiScanHandler)

	// Add generic command execution tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "execute_command",
		Description: "Execute an arbitrary command on the Kali server",
	}, handlers.ExecuteCommandHandler)



	if *httpAddr != "" {
		// Set Gin to release mode for cleaner logs
		gin.SetMode(gin.ReleaseMode)
		ginHandler := gin.New()
		
		// Configure authentication
		authConfig := middleware.NewAuthConfig()
		if authConfig != nil {
			log.Printf("Authentication: Enabled (%s)", authConfig.AuthType)
			ginHandler.Use(middleware.AuthMiddleware(authConfig))
		} else {
			log.Println("Authentication: Disabled (No AUTH_SECRET set)")
			log.Println("WARNING: MCP Server is running without authentication!")
		}

		// Use Gin to wrap the MCP handler
		ginHandler.Any("/*proxyPath", gin.WrapH(mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
			return server
		}, nil)))
		log.Printf("Starting MCP Server with Kali Linux tools and listening at %s", *httpAddr)
		if err := http.ListenAndServe(*httpAddr, ginHandler); err != nil {
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
