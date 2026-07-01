package handler

import (
	"github.com/PritomKarmokar/zipurl/cmd/config"
	"github.com/PritomKarmokar/zipurl/cmd/dts"
	"github.com/PritomKarmokar/zipurl/cmd/model"
	"github.com/PritomKarmokar/zipurl/cmd/repository"
	"github.com/PritomKarmokar/zipurl/cmd/response"
	"github.com/PritomKarmokar/zipurl/cmd/utils"
	"github.com/labstack/echo/v5"
	"github.com/spf13/viper"
	"time"
)

func UrlShortenerHandler(c *echo.Context) error {
	logger := config.GetRequestLogger(c)
	db := config.GetDatabase()

	reqBody := dts.ShortUrlRequest{}
	if err := c.Bind(&reqBody); err != nil {
		logger.Error().Err(err).Msg("failed to bind request body")
		return response.TechnicalError400.ReturnResponse(c, nil)
	}

	if err := c.Validate(reqBody); err != nil {
		logger.Error().Err(err).Msg("Invalid request body")
		return response.DataValidationErr400.ReturnResponse(c, nil)
	}

	id := utils.GenerateULID()
	currentTime := time.Now()
	uniqueToken := utils.EncodeString(id[20:]) // Generating Base62 encoded token from last 7 digits of ulid id

	newUrl := &model.URL{
		ID:          utils.GenerateULID(),
		URL:         reqBody.Url,
		HashedToken: uniqueToken,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
	}
	if err := repository.CreateUrlDBObject(db, newUrl); err != nil {
		logger.Warn().Err(err).Msg("Failed to create url db object")
	}

	logger.Info().Msg("URL shortener DB object created successfully")

	shortUrl := viper.GetString("ZIP_URL_BASE_URL") + "/" + newUrl.HashedToken

	responseData := map[string]interface{}{
		"short_url": shortUrl,
	}
	return response.GenericSuccess200.ReturnResponse(c, responseData)
}
