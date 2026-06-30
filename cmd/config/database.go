package config

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// db connection
var db *gorm.DB

func getDsn() string {
	dbName := viper.GetString("DB_NAME")
	dbUser := viper.GetString("DB_USER")
	dbPassword := viper.GetString("DB_PASSWORD")
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")
	dbSSLMode := viper.GetString("DB_SSLMODE")
	dbTimeZone := viper.GetString("TIME_ZONE")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode, dbTimeZone,
	)
	return dsn
}

func ConnectDB() {
	var err error
	dsn := getDsn()
	slowSQLThreshold := viper.GetDuration("SLOW_SQL_THRESHOLD") * time.Millisecond

	newLogger := logger.New(
		log.New(os.Stdout, "Slow SQL Log: ", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             slowSQLThreshold, // Slow SQL threshold
			LogLevel:                  logger.Warn,      // Log level
			IgnoreRecordNotFoundError: true,             // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,             // Don't include params in the SQL log
			Colorful:                  false,            // Disable color
		},
	)

	zerologLogger := GetLogger()

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		zerologLogger.Fatal().Err(err).Msg("Failed to connect to the database.")
	}
	zerologLogger.Info().Msg("DB connection established successfully.")
}

func GetDatabase() *gorm.DB {
	return db
}
