package main

import (
	"github.com/PritomKarmokar/zipurl/cmd/config"
	"github.com/PritomKarmokar/zipurl/cmd/route"
	"github.com/labstack/echo/v5"
)

func main() {
	e := echo.New()

	config.LoadEnv()
	config.LoggerConfig()
	config.EchoConfig(e)
	route.RegisterRoutes(e)
	config.StartServer(e)
}
