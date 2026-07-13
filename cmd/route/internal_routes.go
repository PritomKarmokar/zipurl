package route

import (
	"github.com/PritomKarmokar/zipurl/cmd/handler"
	"github.com/labstack/echo/v5"
)

func RegisterInternalRoutes(route *echo.Group) {
	route.POST("/url/shorten", handler.UrlShortenerHandler)

	route.POST("/user/signup", handler.UserSignUpHandler)

	route.POST("/user/login", handler.UserLoginHandler)
}
