package main

import (
	"log"

	"github.com/ba0f3/MCP-Kali-Server/pkg/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
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
