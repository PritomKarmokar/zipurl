package route

import (
	"github.com/PritomKarmokar/zipurl/cmd/middleware"
	"github.com/labstack/echo/v5"
)

func RegisterRoutes(e *echo.Echo) {
	// Apply correlation ID tracking globally (first middleware for distributed tracing)
	e.Use(middleware.CorrelationID())

	// base prefix for routes
	basePrefix := e.Group("")

	// Service Routes (Health Checks)
	healthGroup := basePrefix.Group("/health")
	RegisterServiceRoutes(healthGroup)

}
