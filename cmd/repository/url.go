package repository

import (
	"errors"
	"time"

	"github.com/PritomKarmokar/zipurl/cmd/model"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
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

func FetchUrlByToken(db *gorm.DB, token string) (*model.URL, error) {
	start := time.Now()
	var url model.URL

	log.Debug().Str("token", token).Msg("Fetching url db object")

	result := db.
		Where("hashed_token = ?", token).
		Where("expires_at IS NULL OR expires_at > ?", time.Now()).
		Where("max_clicks IS NULL OR click_count < max_clicks").
		First(&url)

	duration := time.Since(start)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Debug().
				Str("operation", "FetchUrlByToken").
				Dur("duration_ms", duration).
				Msg("url db object not found, expired, or click limit reached")
			return nil, nil
		}
		log.Error().
			Err(result.Error).
			Str("operation", "FetchUrlByToken").
			Dur("duration_ms", duration).
			Msg("Failed to fetch url db object")
		return nil, result.Error
	}

	log.Info().
		Str("operation", "FetchUrlByToken").
		Str("hashed_token", token).
		Dur("duration_ms", duration).
		Msg("url db object fetched successfully")
	return &url, nil
}

func FindExistingURLForAnonymousUser(db *gorm.DB, originalUrl string) (*model.URL, error) {
	start := time.Now()
	var url model.URL

	log.Debug().
		Str("url", originalUrl).
		Msg("Fetching URL db object for annoying user")

	result := db.
		Where("url = ?", originalUrl).
		Where("user_id IS NULL OR user_id = ''").
		First(&url)

	duration := time.Since(start)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Debug().
				Str("operation", "FindExistingURLForAnonymousUser").
				Dur("duration_ms", duration).
				Msg("url object not found")
			return nil, nil
		}
		log.Error().
			Err(result.Error).
			Str("operation", "FindExistingURLForAnonymousUser").
			Dur("duration_ms", duration).
			Msg("Failed to fetch url db object")
		return nil, result.Error
	}

	log.Info().
		Str("operation", "FindExistingURLForAnonymousUser").
		Str("url", originalUrl).
		Dur("duration_ms", duration).
		Msg("url db object fetched successful")
	return &url, nil
}

func FindExistingUrlForLoggedInUser(db *gorm.DB, originalUrl string, userId string, currentTime time.Time) (*model.URL, error) {
	start := time.Now()
	var url model.URL

	log.Debug().
		Str("url", originalUrl).
		Msg("Fetching URL db object for logged in User")

	result := db.
		Where("url = ?", originalUrl).
		Where("user_id = ?", userId).
		Where("expires_at > ?", currentTime).
		Where("click_count < max_clicks").
		First(&url)

	duration := time.Since(start)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Debug().
				Str("operation", "FindExistingUrlForLoggedInUser").
				Dur("duration_ms", duration).
				Msg("no reusable url object found")
			return nil, nil
		}
		log.Error().
			Err(result.Error).
			Str("operation", "FindExistingUrlForLoggedInUser").
			Dur("duration_ms", duration).
			Msg("Failed to fetch url db object")
		return nil, result.Error
	}

	log.Info().
		Str("operation", "FindExistingUrlForLoggedInUser").
		Str("url", originalUrl).
		Dur("duration_ms", duration).
		Msg("reusable url db object fetched successfully")

	return &url, nil
}

func IncrementClickCount(db *gorm.DB, token string) error {
	start := time.Now()

	result := db.Model(&model.URL{}).
		Where("hashed_token = ?", token).
		Updates(map[string]interface{}{
			"click_count": gorm.Expr("click_count + 1"),
			"updated_at":  time.Now(),
		})

	duration := time.Since(start)
	if result.Error != nil {
		log.Error().
			Err(result.Error).
			Str("operation", "IncrementClickCount").
			Dur("duration_ms", duration).
			Msg("Failed to increment click count")
		return result.Error
	}
	return nil
}
