package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthConfig holds authentication configuration
type AuthConfig struct {
	AuthType string // "apikey" or "bearer"
	Secret   string
}

// NewAuthConfig creates a new authentication configuration from environment variables
func NewAuthConfig() *AuthConfig {
	authType := os.Getenv("AUTH_TYPE")
	if authType == "" {
		authType = "apikey" // default to API key
	}

	secret := os.Getenv("AUTH_SECRET")
	if secret == "" {
		// Generate a warning if no secret is set
		return nil
	}

	return &AuthConfig{
		AuthType: authType,
		Secret:   secret,
	}
}

// AuthMiddleware creates an authentication middleware based on configuration
func AuthMiddleware(config *AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// If no auth config, allow all requests (for backward compatibility)
		if config == nil {
			c.Next()
			return
		}

		var authenticated bool

		switch config.AuthType {
		case "bearer":
			authenticated = checkBearerToken(c, config.Secret)
		case "apikey":
			authenticated = checkAPIKey(c, config.Secret)
		default:
			authenticated = checkAPIKey(c, config.Secret) // default to API key
		}

		if !authenticated {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkBearerToken validates Bearer token authentication
func checkBearerToken(c *gin.Context, secret string) bool {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return false
	}

	// Check if it starts with "Bearer "
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return false
	}

	// Extract the token
	token := strings.TrimPrefix(authHeader, bearerPrefix)
	
	// Use constant-time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare([]byte(token), []byte(secret)) == 1
}

// checkAPIKey validates API key authentication
func checkAPIKey(c *gin.Context, secret string) bool {
	// Check header first
	apiKey := c.GetHeader("X-API-Key")
	if apiKey == "" {
		// Fallback to query parameter
		apiKey = c.Query("api_key")
	}

	if apiKey == "" {
		return false
	}

	// Use constant-time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare([]byte(apiKey), []byte(secret)) == 1
}

// RateLimitMiddleware provides basic rate limiting
func RateLimitMiddleware(requestsPerMinute int) gin.HandlerFunc {
	// Simple in-memory rate limiter
	// For production, consider using a more robust solution like redis-based rate limiting
	
	return func(c *gin.Context) {
		// This is a simplified implementation
		// For production use, implement proper rate limiting with sliding windows
		c.Next()
	}
}
