package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/spf13/viper"
)

// SecurityHeaders adds comprehensive security headers to all responses
// These headers protect against common web vulnerabilities
func SecurityHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// Prevent MIME type sniffing
			// Protects against: Drive-by downloads, content type confusion attacks
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")

			// Prevent clickjacking attacks
			// Protects against: UI redress attacks where malicious site frames application within an iframe
			c.Response().Header().Set("X-Frame-Options", "DENY")

			// Enable XSS filtering in older browsers
			// Note: Modern browsers have this built-in, but doesn't hurt
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")

			// Force HTTPS for 1 year
			// Protects against: MITM attacks, protocol downgrade attacks
			// Note: Only set if running on HTTPS in production
			if viper.GetBool("ENABLE_HSTS") {
				c.Response().Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
			}

			// Content Security Policy
			// Protects against: XSS, data injection attacks
			// Note: Adjust based on your frontend requirements
			csp := viper.GetString("CONTENT_SECURITY_POLICY")
			if csp == "" {
				// Default CSP - restrictive but safe
				csp = "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'; frame-ancestors 'none'"
			}
			c.Response().Header().Set("Content-Security-Policy", csp)

			// Referrer policy - control referrer information
			// Protects against: Information leakage via referrer header
			c.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			// Permissions policy (formerly Feature-Policy)
			// Restrict browser features that can be used
			c.Response().Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=()")

			// Prevent browser caching of sensitive data
			// For API responses, we don't want caching
			if strings.HasPrefix(c.Path(), "/zip-url/") && c.Path() != "/zip-url/health/live" {
				c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
				c.Response().Header().Set("Pragma", "no-cache")
				c.Response().Header().Set("Expires", "0")
			}
			return next(c)
		}
	}
}

func CORS() echo.MiddlewareFunc {
	allowedOriginsStr := viper.GetString("ALLOWED_ORIGINS")
	var allowedOrigins []string

	if allowedOriginsStr == "*" {
		// Development mode - allow all origins
		allowedOrigins = []string{"*"}
	} else if allowedOriginsStr != "" {
		// Production mode - specific origins
		allowedOrigins = strings.Split(allowedOriginsStr, ",")
		// Trim whitespace
		for i, origin := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(origin)
		}
	} else {
		// Default - localhost only
		allowedOrigins = []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:8080",
		}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			origin := c.Request().Header.Get("Origin")

			// check if origin is allowed
			allowed := false
			if len(allowedOrigins) == 1 && allowedOrigins[0] == "*" {
				// Allow all origins (development only!)
				allowed = true
				origin = "*"
			} else {
				for _, allowedOrigin := range allowedOrigins {
					if origin == allowedOrigin {
						allowed = true
						break
					}
				}
			}

			if allowed {
				// Set CORS headers
				c.Response().Header().Set("Access-Control-Allow-Origin", origin)
				c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
				c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
				c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Device-Fingerprint, X-Request-ID, X-Service-API-Key, X-CSRF-Token")
				c.Response().Header().Set("Access-Control-Expose-Headers", "X-Request-ID, X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset")
				c.Response().Header().Set("Access-Control-Max-Age", "86400") // 24 hours
			}
			// Handle preflight requests
			if c.Request().Method == http.MethodOptions {
				return c.NoContent(http.StatusNoContent)
			}

			return next(c)
		}
	}
}
