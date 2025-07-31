# MCP-Kali-Server

MCP-Kali-Server is a Go-based server application that integrates various Kali Linux tools using the Model Context Protocol (MCP). It provides both HTTP and MCP-based interfaces to execute security tools used for penetration testing and security scanning.

**Note**: This project is a Go fork of the original Python implementation at [https://github.com/Wh0am123/MCP-Kali-Server](https://github.com/Wh0am123/MCP-Kali-Server).

## Features

- **MCP and HTTP Server**: Run the server in MCP mode using an MCP server or in HTTP mode using Gin.
- **Integrated Tools**:
  - Nmap
  - Gobuster
  - Dirb
  - Nikto
  - SQLmap
  - Hydra
  - John the Ripper
  - WPScan
  - Enum4linux

## Project Structure

- `cmd/`: Contains the main executables for the Kali and MCP servers.
- `pkg/tools/`: Implements the tool execution logic.
- `pkg/handlers/`: Defines handlers for tools and commands.
- `go.mod` and `go.sum`: Manage Go dependencies.

## Getting Started

### Prerequisites

- Go 1.24.5 or later
- MCP Go SDK
- Kali Linux tools installed

### Installation

```bash
git clone https://github.com/ba0f3/MCP-Kali-Server.git
cd MCP-Kali-Server
go mod tidy
```

### Running the Server

#### MCP Mode

```bash
go run ./cmd/mcp-server -debug
```

#### HTTP Mode

```bash
SERVER_MODE=gin go run ./cmd/kali-server
```

The server starts on port 5000 by default and hosts multiple endpoints for tool operations.

## Authentication

The HTTP server supports authentication to protect your endpoints. Authentication is configured through environment variables:

- `AUTH_SECRET`: The secret key/token for authentication (required to enable auth)
- `AUTH_TYPE`: The authentication method - either `apikey` (default) or `bearer`

### API Key Authentication (default)
Set the secret and include it in requests:
```bash
# Set the secret
export AUTH_SECRET=your-secret-key

# Include in header
curl -H "X-API-Key: your-secret-key" http://localhost:5000/api/tools/nmap -d '{...}'

# Or as query parameter
curl http://localhost:5000/api/tools/nmap?api_key=your-secret-key -d '{...}'
```

### Bearer Token Authentication
```bash
# Set auth type and secret
export AUTH_TYPE=bearer
export AUTH_SECRET=your-bearer-token

# Include in Authorization header
curl -H "Authorization: Bearer your-bearer-token" http://localhost:5000/api/tools/nmap -d '{...}'
```

## Usage

### Example Commands

- Nmap scan:
  ```bash
  curl -X POST http://localhost:5000/api/tools/nmap -d '{"target": "example.com", "scan_type": "-sS"}'
  ```

- WPScan analysis:
  ```bash
  curl -X POST http://localhost:5000/api/tools/wpscan -d '{"url": "http://example.com"}'
  ```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
