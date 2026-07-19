package config

import (
	"sync"

	"github.com/PritomKarmokar/zipurl/cmd/service"
	"github.com/spf13/viper"
)

type Clients struct {
	TokenClient *service.TokenClient
	RedisClient *service.RedisClient
}

var (
	ClientsRef      Clients
	redisAvailable  bool
	redisAvailMutex sync.RWMutex
)

func IsRedisAvailable() bool {
	redisAvailMutex.RLock()
	defer redisAvailMutex.RUnlock()
	return redisAvailable
}

func setRedisAvailable(available bool) {
	redisAvailMutex.Lock()
	defer redisAvailMutex.Unlock()
	redisAvailable = available
}

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

	var redisClient *service.RedisClient
	redisClient, err = service.NewRedisClient(viper.GetString("REDIS_DSN"))
	if err != nil {
		logger.Error().Err(err).Msg("Failed to Initiate Redis client")
		setRedisAvailable(false)
		redisClient = nil
	} else {
		logger.Info().Msg("Redis Connected Successfully")
		setRedisAvailable(true)
	}

	ClientsRef.TokenClient = tokenClient
	ClientsRef.RedisClient = redisClient

	logger.Info().Msg("All clients initialized successfully")
}

func GetClients() *Clients {
	return &ClientsRef
}
