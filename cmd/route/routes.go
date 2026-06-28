package route

import "github.com/labstack/echo/v5"

func RegisterRoutes(e *echo.Echo) {
	// base prefix for routes
	basePrefix := e.Group("")

	// Service Routes (Health Checks)
	healthGroup := basePrefix.Group("/health")
	RegisterServiceRoutes(healthGroup)
	
}
