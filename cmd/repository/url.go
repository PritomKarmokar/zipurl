package repository

import (
	"errors"
	"github.com/PritomKarmokar/zipurl/cmd/model"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

func CreateUrlDBObject(db *gorm.DB, data *model.URL) error {
	start := time.Now()

	result := db.Create(data)
	duration := time.Since(start)

	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Dur("duration", duration).
			Msg("Failed to create url db object")
		return result.Error
	}

	log.Debug().
		//Str("url_db_object_id", result.ID).
		Dur("duration_ms", duration).
		Msg("URL DB object created")

	return nil
}

func FetchUrlDBObject(db *gorm.DB, token string) (*model.URL, error) {
	start := time.Now()
	var url model.URL

	log.Debug().
		Str("token", token).
		Msg("Fetching url db object")

	result := db.Where("hashed_token = ?", token).First(&url)

	duration := time.Since(start)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Debug().
				Str("operation", "FetchUrlDBObject").
				Dur("duration_ms", duration).
				Msg("url db object not found")
			return nil, nil
		}
		log.Error().
			Err(result.Error).
			Str("operation", "FetchUrlDBObject").
			Dur("duration_ms", duration).
			Msg("Failed to fetch url db object")
		return nil, result.Error
	}

	log.Info().
		Str("operation", "FetchUrlDBObject").
		Str("hashed_token", token).
		Dur("duration_ms", duration).
		Msg("url db object fetched successfully")

	return &url, nil
}
