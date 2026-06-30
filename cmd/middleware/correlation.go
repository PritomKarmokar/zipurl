package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func CorrelationID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// Check if correlation ID already exists (from upstream service)
			correlationID := c.Request().Header.Get("X-Correlation-ID")

			// Generate new correlation ID if not present
			if correlationID == "" {
				correlationID = uuid.New().String()
			}

			// Also handle X-Request-ID for backward compatibility
			requestID := c.Request().Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
			}

			// Set in response headers for client tracking
			c.Response().Header().Set("X-Correlation-ID", correlationID)
			c.Response().Header().Set("X-Request-ID", requestID)

			// Store in context for use in handlers and logging
			c.Set("correlation_id", correlationID)
			c.Set("request_id", requestID)

			return next(c)
		}
	}
}
