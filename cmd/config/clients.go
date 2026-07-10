package config

import (
	"github.com/PritomKarmokar/zipurl/cmd/service"
	"github.com/spf13/viper"
)

type Clients struct {
	TokenClient *service.TokenClient
}

var (
	ClientsRef Clients
)

func LoadClients() {
	logger := GetLogger()

	// Token Client
	tokenClient, err := service.NewTokenClient(
		viper.GetString("JWT_VERIFYING_KEY"),
		viper.GetString("JWT_SIGNING_KEY"),
		viper.GetString("JWT_ALGORITHM"),
		viper.GetStringSlice("JWT_AUDIENCE"),
		viper.GetString("JWT_ISSUER"),
		viper.GetInt64("ACCESS_TOKEN_EXPIRY"),
		viper.GetInt64("REFRESH_TOKEN_EXPIRY"),
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initiate token client")
	}
	ClientsRef.TokenClient = tokenClient

	logger.Info().Msg("All clients initialized successfully")
}

func GetClients() *Clients {
	return &ClientsRef
}
