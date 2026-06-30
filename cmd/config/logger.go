package config

import (
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

func GetLogger() *zerolog.Logger {
	return &log.Logger
}
