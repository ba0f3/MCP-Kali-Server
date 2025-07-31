package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Wh0am123/MCP-Kali-Server/pkg/executor"
	"github.com/gin-gonic/gin"
)

func genericCommand(c *gin.Context) {
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

	result, err := executor.ExecuteCommand(command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Helper function to safely get string from map
func getStringParam(data map[string]interface{}, key string, defaultValue string) string {
	if val, ok := data[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return defaultValue
}

func nmap(c *gin.Context) {
	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
		return
	}

	target := getStringParam(data, "target", "")
	if target == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Target parameter is required"})
		return
	}

	scanType := getStringParam(data, "scan_type", "-sCV")
	ports := getStringParam(data, "ports", "")
	additionalArgs := getStringParam(data, "additional_args", "-T4 -Pn")

	command := fmt.Sprintf("nmap %s", scanType)
	if ports != "" {
		command += fmt.Sprintf(" -p %s", ports)
	}
	if additionalArgs != "" {
		command += fmt.Sprintf(" %s", additionalArgs)
	}
	command += fmt.Sprintf(" %s", target)

	result, err := executor.ExecuteCommand(command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}


func gobuster(c *gin.Context) {
	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
		return
	}

	url := getStringParam(data, "url", "")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is required"})
		return
	}

	mode := getStringParam(data, "mode", "dir")
	// Validate mode
	validModes := map[string]bool{"dir": true, "dns": true, "fuzz": true, "vhost": true}
	if !validModes[mode] {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid mode: %s. Must be one of: dir, dns, fuzz, vhost", mode)})
		return
	}

	wordlist := getStringParam(data, "wordlist", "/usr/share/wordlists/dirb/common.txt")
	additionalArgs := getStringParam(data, "additional_args", "")

	command := fmt.Sprintf("gobuster %s -u %s -w %s", mode, url, wordlist)
	if additionalArgs != "" {
		command += fmt.Sprintf(" %s", additionalArgs)
	}

	result, err := executor.ExecuteCommand(command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
func dirb(c *gin.Context) {
	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
		return
	}

	url := getStringParam(data, "url", "")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is required"})
		return
	}

	wordlist := getStringParam(data, "wordlist", "/usr/share/wordlists/dirb/common.txt")
	additionalArgs := getStringParam(data, "additional_args", "")

	command := fmt.Sprintf("dirb %s %s", url, wordlist)
	if additionalArgs != "" {
		command += fmt.Sprintf(" %s", additionalArgs)
	}

	result, err := executor.ExecuteCommand(command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func nikto(c *gin.Context) {
	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
		return
	}

	target := getStringParam(data, "target", "")
	if target == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Target parameter is required"})
		return
	}

	additionalArgs := getStringParam(data, "additional_args", "")

	command := fmt.Sprintf("nikto -h %s", target)
	if additionalArgs != "" {
		command += fmt.Sprintf(" %s", additionalArgs)
	}

	result, err := executor.ExecuteCommand(command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func sqlmap(c *gin.Context) {
	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
		return
	}

	url := getStringParam(data, "url", "")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is required"})
		return
	}

	postData := getStringParam(data, "data", "")
	additionalArgs := getStringParam(data, "additional_args", "")

	command := fmt.Sprintf("sqlmap -u %s --batch", url)
	if postData != "" {
		command += fmt.Sprintf(" --data=\"%s\"", postData)
	}
	if additionalArgs != "" {
		command += fmt.Sprintf(" %s", additionalArgs)
	}

	result, err := executor.ExecuteCommand(command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func metasploit(c *gin.Context) {
	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
		return
	}

	module := getStringParam(data, "module", "")
	if module == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Module parameter is required"})
		return
	}

	options := data["options"].(map[string]interface{})

	// Format options
	var optionsStr strings.Builder
	for key, value := range options {
		optionsStr.WriteString(fmt.Sprintf(" set %s %v\n", key, value))
	}

	resourceContent := fmt.Sprintf("use %s\n%s exploit\n", module, optionsStr.String())
	tempFile, err := os.CreateTemp("", "mcp_msf_resource.rc")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary file."})
		return
	}
	defer os.Remove(tempFile.Name())

	tempFile.WriteString(resourceContent)
	tempFile.Close()

	command := fmt.Sprintf("msfconsole -q -r %s", tempFile.Name())
	result, err := executor.ExecuteCommand(command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
func hydra(c *gin.Context) {
	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
		return
	}

	target := getStringParam(data, "target", "")
	if target == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Target parameter is required"})
		return
	}

	service := getStringParam(data, "service", "")
	if service == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service parameter is required"})
		return
	}

	username := getStringParam(data, "username", "")
	usernameFile := getStringParam(data, "username_file", "")
	password := getStringParam(data, "password", "")
	passwordFile := getStringParam(data, "password_file", "")

	if username == "" && usernameFile == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or username_file parameter is required"})
		return
	}

	if password == "" && passwordFile == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password or password_file parameter is required"})
		return
	}

	additionalArgs := getStringParam(data, "additional_args", "")

	command := fmt.Sprintf("hydra -t 4")

	if username != "" {
		command += fmt.Sprintf(" -l %s", username)
	} else {
		command += fmt.Sprintf(" -L %s", usernameFile)
	}

	if password != "" {
		command += fmt.Sprintf(" -p %s", password)
	} else {
		command += fmt.Sprintf(" -P %s", passwordFile)
	}

	if additionalArgs != "" {
		command += fmt.Sprintf(" %s", additionalArgs)
	}

	command += fmt.Sprintf(" %s %s", target, service)

	result, err := executor.ExecuteCommand(command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func john(c *gin.Context) {
	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
		return
	}

	hashFile := getStringParam(data, "hash_file", "")
	if hashFile == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hash file parameter is required"})
		return
	}

	wordlist := getStringParam(data, "wordlist", "/usr/share/wordlists/rockyou.txt")
	formatType := getStringParam(data, "format", "")
	additionalArgs := getStringParam(data, "additional_args", "")

	command := fmt.Sprintf("john")

	if formatType != "" {
		command += fmt.Sprintf(" --format=%s", formatType)
	}

	if wordlist != "" {
		command += fmt.Sprintf(" --wordlist=%s", wordlist)
	}

	if additionalArgs != "" {
		command += fmt.Sprintf(" %s", additionalArgs)
	}

	command += fmt.Sprintf(" %s", hashFile)
	result, err := executor.ExecuteCommand(command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func wpscan(c *gin.Context) {
	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
		return
	}

	url := getStringParam(data, "url", "")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is required"})
		return
	}

	additionalArgs := getStringParam(data, "additional_args", "")

	command := fmt.Sprintf("wpscan --url %s", url)
	if additionalArgs != "" {
		command += fmt.Sprintf(" %s", additionalArgs)
	}

	result, err := executor.ExecuteCommand(command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func enum4linux(c *gin.Context) {
	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
		return
	}

	target := getStringParam(data, "target", "")
	if target == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Target parameter is required"})
		return
	}

	additionalArgs := getStringParam(data, "additional_args", "-a")

	command := fmt.Sprintf("enum4linux %s %s", additionalArgs, target)
	result, err := executor.ExecuteCommand(command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func healthCheck(c *gin.Context) {

essentialTools := []string{"nmap", "gobuster", "dirb", "nikto"}

toolsStatus := map[string]bool{}

	for _, tool := range essentialTools {
		command := fmt.Sprintf("which %s", tool)
		result, _ := executor.ExecuteCommand(command)
		success := strings.TrimSpace(result.Stdout) != ""
		toolsStatus[tool] = success
	}

	allEssentialToolsAvailable := true
	for _, available := range toolsStatus {
		if !available {
			allEssentialToolsAvailable = false
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"message": "Kali Linux Tools API Server is running",
		"tools_status": toolsStatus,
		"all_essential_tools_available": allEssentialToolsAvailable,
	})
}

func main() {
	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Setup routes
	r.POST("/api/command", genericCommand)
	r.POST("/api/tools/nmap", nmap)
	r.POST("/api/tools/gobuster", gobuster)
	r.POST("/api/tools/dirb", dirb)
	r.POST("/api/tools/nikto", nikto)
	r.POST("/api/tools/sqlmap", sqlmap)
	r.POST("/api/tools/metasploit", metasploit)
	r.POST("/api/tools/hydra", hydra)
	r.POST("/api/tools/john", john)
	r.POST("/api/tools/wpscan", wpscan)
	r.POST("/api/tools/enum4linux", enum4linux)
	r.GET("/health", healthCheck)

	// Start server
	log.Printf("Starting server on port 5000")
	if err := r.Run(":5000"); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
