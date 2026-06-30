package config

import "github.com/spf13/viper"

func LoadEnv() {
	logger := GetLogger()

	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		viper.AutomaticEnv()
		logger.Info().Err(err).Msg("Failed to load config file. ENV loaded from AutomaticEnv()")
	} else {
		logger.Info().Msg("ENV loaded from .env")
	}
}
