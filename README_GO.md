# MCP Kali Server - Go Version

This is a Go implementation of the MCP Kali Server, converted from the original Python version.

## Components

### 1. Kali API Server (`cmd/kali-server`)
A Gin-based HTTP API server that runs on the Kali Linux machine and executes security tools.

### 2. MCP Client (To be implemented)
An MCP client that connects to the Kali API server and exposes tools via the MCP protocol.

## Project Structure

```
MCP-Kali-Server/
├── cmd/
│   ├── kali-server/     # Kali Linux API server
│   │   └── main.go
│   └── mcp-server/       # MCP server (client to Kali API)
│       └── main.go
├── pkg/
│   ├── executor/         # Command execution package
│   │   └── executor.go
│   └── kalitools/        # Kali tools client package
│       └── client.go
├── go.mod
├── go.sum
└── README_GO.md
```

## Building

### Build Kali API Server
```bash
cd cmd/kali-server
go build -o kali-server main.go
```

### Build MCP Server
```bash
cd cmd/mcp-server
go build -o mcp-server main.go
```

## Running

### 1. Start the Kali API Server (on Linux)
```bash
./kali-server
```

The server will start on port 5000 by default.

### 2. Start the MCP Server (on any machine)
```bash
./mcp-server --server http://LINUX_IP:5000
```

## API Endpoints

The Kali API server provides the following endpoints:

- `POST /api/command` - Execute arbitrary command
- `POST /api/tools/nmap` - Execute Nmap scan
- `POST /api/tools/gobuster` - Execute Gobuster scan
- `POST /api/tools/dirb` - Execute Dirb scan
- `POST /api/tools/nikto` - Execute Nikto scan
- `POST /api/tools/sqlmap` - Execute SQLmap scan
- `POST /api/tools/metasploit` - Execute Metasploit module
- `POST /api/tools/hydra` - Execute Hydra attack
- `POST /api/tools/john` - Execute John the Ripper
- `POST /api/tools/wpscan` - Execute WPScan
- `POST /api/tools/enum4linux` - Execute Enum4linux
- `GET /health` - Health check endpoint

## Features

- **Command Execution with Timeout**: All commands execute with a configurable timeout (default 3 minutes)
- **Partial Results**: If a command times out but produces output, the partial results are returned
- **Error Handling**: Comprehensive error handling and logging
- **Health Checks**: Built-in health check endpoint that verifies essential tools are available

## Configuration

### Environment Variables
- `API_PORT` - Port for the API server (default: 5000)

### Command Line Flags

#### kali-server
```bash
./kali-server
```

#### mcp-server
```bash
./mcp-server --server http://localhost:5000 --timeout 5m --debug
```

Flags:
- `--server` - Kali API server URL (default: http://localhost:5000)
- `--timeout` - Request timeout (default: 5m)
- `--debug` - Enable debug logging

## Security Considerations

⚠️ **WARNING**: This tool allows remote command execution. Use with caution and ensure:
- The API server is not exposed to untrusted networks
- Use authentication/authorization in production
- Validate and sanitize all inputs
- Run with minimal required privileges

## Development

### Adding New Tools

To add a new tool to the Kali API server:

1. Create a new handler function in `cmd/kali-server/main.go`
2. Add the route in the `main()` function
3. Update the MCP server to expose the new tool

Example:
```go
func newTool(c *gin.Context) {
    var data map[string]interface{}
    if err := c.BindJSON(&data); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
        return
    }

    // Extract parameters
    target := getStringParam(data, "target", "")
    
    // Build command
    command := fmt.Sprintf("newtool %s", target)
    
    // Execute
    result, err := executor.ExecuteCommand(command)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, result)
}
```

## License

Same as the original Python version - see LICENSE file.
