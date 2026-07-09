package repository

import (
	"time"

	"github.com/PritomKarmokar/zipurl/cmd/model"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func CreateNewUserObject(db *gorm.DB, user *model.User) error {
	start := time.Now()

	result := db.Create(user)
	duration := time.Since(start)

	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Dur("duration", duration).
			Msg("Failed to create a new user")
		return result.Error
	}

	log.Debug().
		Dur("duration_ms", duration).
		Msg("New User signup successful")

	return nil
}
