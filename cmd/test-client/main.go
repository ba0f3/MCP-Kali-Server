package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	var (
		serverURL = flag.String("server", "http://localhost:5000", "Kali API server URL")
		command   = flag.String("cmd", "whoami", "Command to execute")
		endpoint  = flag.String("endpoint", "command", "API endpoint (command, health)")
	)
	flag.Parse()

	switch *endpoint {
	case "command":
		testCommand(*serverURL, *command)
	case "health":
		testHealth(*serverURL)
	case "nmap":
		testNmap(*serverURL)
	default:
		log.Fatalf("Unknown endpoint: %s", *endpoint)
	}
}

func testCommand(serverURL, command string) {
	data := map[string]string{
		"command": command,
	}
	
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(serverURL+"/api/command", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: %s\n", string(body))
}

func testHealth(serverURL string) {
	resp, err := http.Get(serverURL + "/health")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Health Response: %s\n", string(body))
}

func testNmap(serverURL string) {
	data := map[string]interface{}{
		"target": "127.0.0.1",
		"scan_type": "-sV",
		"ports": "22,80,443",
	}
	
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(serverURL+"/api/tools/nmap", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Nmap Response: %s\n", string(body))
}
