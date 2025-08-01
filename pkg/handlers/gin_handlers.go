package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/ba0f3/MCP-Kali-Server/pkg/executor"
	"github.com/gin-gonic/gin"
)

// Helper function to safely get string from map
func getStringParam(data map[string]interface{}, key string, defaultValue string) string {
	if val, ok := data[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return defaultValue
}

func GenericCommandHandler(c *gin.Context) {
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

func NmapHandler(c *gin.Context) {
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

func GobusterHandler(c *gin.Context) {
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

func DirbHandler(c *gin.Context) {
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

func NiktoHandler(c *gin.Context) {
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

func SqlmapHandler(c *gin.Context) {
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

func MetasploitHandler(c *gin.Context) {
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

func HydraHandler(c *gin.Context) {
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

func JohnHandler(c *gin.Context) {
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

func WpscanHandler(c *gin.Context) {
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

func Enum4linuxHandler(c *gin.Context) {
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

func Sublist3rHandler(c *gin.Context) {
	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
		return
	}

	domain := getStringParam(data, "domain", "")
	if domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Domain parameter is required"})
		return
	}

	// Safely get boolean values
	bruteForce := false
	if val, ok := data["bruteforce"]; ok {
		if bVal, ok := val.(bool); ok {
			bruteForce = bVal
		}
	}

	ports := getStringParam(data, "ports", "")
	
	// Safely get threads value
	threads := 0
	if val, ok := data["threads"]; ok {
		if fVal, ok := val.(float64); ok {
			threads = int(fVal)
		}
	}
	
	engines := getStringParam(data, "engines", "")
	
	// Safely get verbose value
	verbose := false
	if val, ok := data["verbose"]; ok {
		if bVal, ok := val.(bool); ok {
			verbose = bVal
		}
	}
	additionalArgs := getStringParam(data, "additional_args", "")

	command := fmt.Sprintf("sublist3r -d %s", domain)
	if bruteForce {
		command += " -b"
	}
	if ports != "" {
		command += fmt.Sprintf(" -p %s", ports)
	}
	if threads > 0 {
		command += fmt.Sprintf(" -t %d", threads)
	} else {
		command += " -t 10"
	}
	if engines != "" {
		command += fmt.Sprintf(" -e %s", engines)
	}
	if verbose {
		command += " -v"
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

func HealthCheckHandler(c *gin.Context) {

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
