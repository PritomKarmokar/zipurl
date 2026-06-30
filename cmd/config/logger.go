package config

import (
	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func setLevel(level string) {
	var l zerolog.Level
	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	case "trace":
		l = zerolog.TraceLevel
	default:
		l = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(l)
}

func LoggerConfig() {
	logLevel := viper.GetString("LOG_LEVEL_API")
	setLevel(logLevel)
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
}

func GetRequestLogger(c *echo.Context) zerolog.Logger {
	logger := log.Logger.With().
		Str("x-request-id", c.Response().Header().Get("X-Request-ID")).
		Str("x-correlation-id", c.Response().Header().Get("X-Correlation-ID")).
		Logger()

	// Add correlation_id and request_id from context if available
	if correlationID, ok := c.Get("correlation_id").(string); ok && correlationID != "" {
		logger = logger.With().Str("correlation-id", correlationID).Logger()
	}
	if requestID, ok := c.Get("request_id").(string); ok && requestID != "" {
		logger = logger.With().Str("request_id", requestID).Logger()
	}

	return logger
}

func GetLogger() *zerolog.Logger {
	return &log.Logger
}
