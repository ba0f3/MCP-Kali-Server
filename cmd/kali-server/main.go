package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ba0f3/MCP-Kali-Server/pkg/executor"
	"github.com/ba0f3/MCP-Kali-Server/pkg/handlers"
	"github.com/ba0f3/MCP-Kali-Server/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Define command-line flags
	timeout := flag.Int("timeout", 900, "Command execution timeout in seconds")
	port := flag.Int("port", 5000, "Port to listen on")
	flag.Parse()

	// Set the global command timeout
	executor.SetGlobalTimeout(time.Duration(*timeout) * time.Second)

	// Determine mode of the server from environment variable or configuration
	mode := os.Getenv("SERVER_MODE")
	if mode == "" {
		mode = "gin" // default to Gin HTTP mode
	}

	// Print server configuration
	log.Println("=== MCP-Kali-Server Configuration ===")
	log.Printf("Server Mode: %s", mode)
	log.Printf("Command Timeout: %d seconds", *timeout)
	log.Printf("Port: %d", *port)

	if mode == "mcp" {
		// Create the MCP server
		server := mcp.NewServer(&mcp.Implementation{Name: "Kali Server", Version: "v1.0.0"}, nil)

		// Define tool handlers here
		// Example: mcp.AddTool(server, mcp.Tool{Name: "nmap", Description: "Network exploration tool"}, nmapHandler)

		// Create MCP streamable HTTP handler
		handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return server }, nil)

		// Start the MCP server
		log.Printf("Starting MCP streamable HTTP server on port %d", *port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), handler); err != nil {
			log.Fatalf("Could not start MCP server: %v", err)
		}
	} else {
		// Setup Gin router
		gin.SetMode(gin.ReleaseMode)
		r := gin.Default()

		// Configure authentication
		authConfig := middleware.NewAuthConfig()
		if authConfig != nil {
			log.Printf("Authentication: Enabled (%s)", authConfig.AuthType)
			r.Use(middleware.AuthMiddleware(authConfig))
		} else {
			log.Println("Authentication: Disabled (No AUTH_SECRET set)")
			log.Println("WARNING: Server is running without authentication!")
		}
		log.Println("=====================================")

		// Setup routes
		r.POST("/api/command", handlers.GenericCommandHandler)
		r.POST("/api/stream/command", handlers.StreamCommandHandler)
		r.POST("/api/tools/nmap", handlers.NmapHandler)
		r.POST("/api/tools/gobuster", handlers.GobusterHandler)
		r.POST("/api/tools/dirb", handlers.DirbHandler)
		r.POST("/api/tools/nikto", handlers.NiktoHandler)
		r.POST("/api/tools/sqlmap", handlers.SqlmapHandler)
		r.POST("/api/tools/metasploit", handlers.MetasploitHandler)
		r.POST("/api/tools/hydra", handlers.HydraHandler)
		r.POST("/api/tools/john", handlers.JohnHandler)
		r.POST("/api/tools/wpscan", handlers.WpscanHandler)
		r.POST("/api/tools/enum4linux", handlers.Enum4linuxHandler)
		r.GET("/health", handlers.HealthCheckHandler)

		// Start the Gin server
		log.Printf("Starting Gin HTTP server on port %d", *port)
		if err := r.Run(fmt.Sprintf(":%d", *port)); err != nil {
			log.Fatalf("Could not start Gin server: %v", err)
		}
	}
}
