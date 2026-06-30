package route

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func RegisterServiceRoutes(route *echo.Group) {
	route.GET("/live", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "alive",
		})
	})
}
