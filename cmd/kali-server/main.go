package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Wh0am123/MCP-Kali-Server/pkg/handlers"
	"github.com/gin-gonic/gin"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Determine mode of the server from environment variable or configuration
	mode := os.Getenv("SERVER_MODE")

	if mode == "mcp" {
		// Create the MCP server
		server := mcp.NewServer(&mcp.Implementation{Name: "Kali Server", Version: "v1.0.0"}, nil)

		// Define tool handlers here
		// Example: mcp.AddTool(server, mcp.Tool{Name: "nmap", Description: "Network exploration tool"}, nmapHandler)

		// Create MCP streamable HTTP handler
		handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return server }, nil)

		// Start the MCP server
		log.Printf("Starting MCP streamable HTTP server on port 5000")
		if err := http.ListenAndServe(":5000", handler); err != nil {
			log.Fatalf("Could not start MCP server: %v", err)
		}
	} else {
		// Setup Gin router
		gin.SetMode(gin.ReleaseMode)
		r := gin.Default()

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
		log.Printf("Starting Gin HTTP server on port 5000")
		if err := r.Run(":5000"); err != nil {
			log.Fatalf("Could not start Gin server: %v", err)
		}
	}
}
