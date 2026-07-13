package route

import (
	"github.com/PritomKarmokar/zipurl/cmd/handler"
	"github.com/PritomKarmokar/zipurl/cmd/middleware"
	"github.com/labstack/echo/v5"
)

func RegisterInternalRoutes(route *echo.Group) {
	route.POST("/url/shorten", handler.UrlShortenerHandler)

	route.POST("/user/signup", handler.UserSignUpHandler)

	route.POST("/user/login", handler.UserLoginHandler)
}

func RegisterProtectedInternalRoutes(route *echo.Group) {
	route.Use(middleware.JwtTokenAuth())
}
