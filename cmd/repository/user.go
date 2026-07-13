package repository

import (
	"errors"
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

func UserExistsByEmail(db *gorm.DB, email string) (bool, error) {
	start := time.Now()
	var user model.User

	result := db.Where("email = ?", email).First(&user)
	duration := time.Since(start)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Debug().
				Str("operation", "FindUserWithEmailAddress").
				Dur("duration_ms", duration).
				Msg("User not found with this email")
			return false, nil
		}
		log.Error().
			Err(result.Error).
			Str("email", email).
			Str("operation", "FindUserWithEmailAddress").
			Dur("duration_ms", duration).
			Msg("Unexpected error occurred")
		return false, result.Error
	}

	log.Debug().
		Str("operation", "FindUserWithEmailAddress").
		Str("email", email).
		Dur("duration_ms", duration).
		Msg("User with this email address found")

	return true, nil
}

func GetUserByEmail(db *gorm.DB, email string) (*model.User, error) {
	start := time.Now()
	var user model.User

	log.Debug().
		Str("operation", "GetUserByEmail").
		Str("email", email).
		Msg("Getting user by email")

	result := db.Where("email = ?", email).First(&user)
	duration := time.Since(start)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Debug().
				Str("operation", "GetUserByEmail").
				Str("email", email).
				Dur("duration_ms", duration).
				Msg("User not found with this email")
			return nil, nil
		}
		log.Error().
			Err(result.Error).
			Str("operation", "GetUserByEmail").
			Dur("duration_ms", duration).
			Msg("Failed to get user by email")
		return nil, result.Error
	}

	log.Debug().
		Str("operation", "GetUserByEmail").
		Str("email", email).
		Dur("duration_ms", duration).
		Msg("User found with this email")

	return &user, nil
}
