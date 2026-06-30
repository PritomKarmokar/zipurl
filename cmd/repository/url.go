package repository

import (
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
