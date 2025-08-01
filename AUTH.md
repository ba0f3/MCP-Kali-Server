# Authentication for MCP Server

The MCP server now supports authentication middleware to secure your endpoints.

## Configuration

Authentication is configured using environment variables:

- `AUTH_TYPE`: Type of authentication to use. Options are:
  - `apikey` (default): API key authentication
  - `bearer`: Bearer token authentication
  
- `AUTH_SECRET`: The secret key/token used for authentication

## Usage

### Without Authentication (Development Only)

If no `AUTH_SECRET` is set, the server runs without authentication:

```bash
./mcp-server -http :8080
```

⚠️ **WARNING**: Running without authentication is not recommended for production!

### With API Key Authentication

```bash
export AUTH_TYPE=apikey
export AUTH_SECRET=your-secret-api-key
./mcp-server -http :8080
```

Clients must include the API key in their requests:

```bash
# Using header
curl -H "X-API-Key: your-secret-api-key" http://localhost:8080/...

# Using query parameter
curl http://localhost:8080/...?api_key=your-secret-api-key
```

### With Bearer Token Authentication

```bash
export AUTH_TYPE=bearer
export AUTH_SECRET=your-bearer-token
./mcp-server -http :8080
```

Clients must include the Bearer token in the Authorization header:

```bash
curl -H "Authorization: Bearer your-bearer-token" http://localhost:8080/...
```

## Security Considerations

1. **Use Strong Secrets**: Generate strong, random secrets for production use
2. **HTTPS**: Always use HTTPS in production to prevent token interception
3. **Token Rotation**: Regularly rotate your authentication secrets
4. **Constant-Time Comparison**: The middleware uses constant-time comparison to prevent timing attacks

## Example: Running with systemd (Linux)

Create a service file with authentication:

```ini
[Unit]
Description=MCP Kali Server
After=network.target

[Service]
Type=simple
User=kali
Environment="AUTH_TYPE=apikey"
Environment="AUTH_SECRET=your-secure-api-key"
ExecStart=/usr/local/bin/mcp-server -http :8080
Restart=always

[Install]
WantedBy=multi-user.target
```

## Example: Running with Windows Service

When installing as a Windows service with authentication:

```powershell
# Set system environment variables
[System.Environment]::SetEnvironmentVariable("AUTH_TYPE", "apikey", "Machine")
[System.Environment]::SetEnvironmentVariable("AUTH_SECRET", "your-secure-api-key", "Machine")

# Install the service
./mcp-server -install-service -service-name mcp-kali-server -service-port :8080
```

## Testing Authentication

Test that authentication is working:

```bash
# Should return 401 Unauthorized
curl -v http://localhost:8080/test

# Should pass authentication (but may return 400 due to MCP protocol requirements)
curl -v -H "X-API-Key: your-secret-api-key" http://localhost:8080/test
```
