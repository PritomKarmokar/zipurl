package config

import (
	"context"
	"github.com/labstack/echo/v5"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartServer(e *echo.Echo) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sc := echo.StartConfig{
		Address:         ":" + viper.GetString("SERVER_PORT"),
		GracefulTimeout: 5 * time.Second,
	}

	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("Failed to start server", "error", err)
	}
}
