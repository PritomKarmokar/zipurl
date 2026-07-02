package route

import (
	"github.com/PritomKarmokar/zipurl/cmd/middleware"
	"github.com/labstack/echo/v5"
)

func RegisterRoutes(e *echo.Echo) {
	// Apply correlation ID tracking globally (first middleware for distributed tracing)
	e.Use(middleware.CorrelationID())

	// Apply security headers globally
	e.Use(middleware.SecurityHeaders())

	// Apply CORS for web (configure allowed origins in env)
	e.Use(middleware.CORS())

	// base prefix for routes
	basePrefix := e.Group("/zip-url")

	// Service Routes (Health Checks)
	healthGroup := basePrefix.Group("/health")
	RegisterServiceRoutes(healthGroup)

	internalRoute := basePrefix.Group("/api/v1")
	RegisterInternalRoutes(internalRoute)

	// Redirect Routes (no base prefix provided)
	redirectRoute := e.Group("/")
	RegisterRedirectRoutes(redirectRoute)
}
