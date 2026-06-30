package config

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func EchoConfig(e *echo.Echo) {
	
	e.Use(middleware.RequestID())
	e.Use(middleware.BodyLimit(2_097_152)) // 2 MB

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
}
