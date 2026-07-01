package route

import (
	"github.com/PritomKarmokar/zipurl/cmd/handler"
	"github.com/labstack/echo/v5"
)

func RegisterRedirectRoutes(route *echo.Group) {
	route.GET(":token", handler.UrlRedirectHandler)
}
