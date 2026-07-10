package handler

import (
	"net/http"
	"time"

	"github.com/PritomKarmokar/zipurl/cmd/config"
	"github.com/PritomKarmokar/zipurl/cmd/dts"
	"github.com/PritomKarmokar/zipurl/cmd/model"
	"github.com/PritomKarmokar/zipurl/cmd/repository"
	"github.com/PritomKarmokar/zipurl/cmd/response"
	"github.com/PritomKarmokar/zipurl/cmd/utils"
	"github.com/labstack/echo/v5"
	"github.com/spf13/viper"
)

func UrlShortenerHandler(c *echo.Context) error {
	logger := config.GetRequestLogger(c)
	db := config.GetDatabase()

	reqBody := dts.ShortUrlRequest{}
	if err := c.Bind(&reqBody); err != nil {
		logger.Error().Err(err).Msg("failed to bind request body")
		return response.TechnicalError.ReturnResponse(c, nil)
	}

	if err := c.Validate(reqBody); err != nil {
		logger.Error().Err(err).Msg("Invalid request body")
		return response.DataValidationErr.ReturnResponse(c, nil)
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
		return response.ShortURLCreationFailed.ReturnResponse(c, nil)
	}

	logger.Info().Msg("URL Shortener DB object created successfully")

	shortUrl := viper.GetString("ZIP_URL_BASE_URL") + "/" + newUrl.HashedToken

	responseData := map[string]interface{}{
		"short_url": shortUrl,
	}
	return response.GenericSuccess200.ReturnResponse(c, responseData)
}

func UrlRedirectHandler(c *echo.Context) error {
	logger := config.GetRequestLogger(c)
	db := config.GetDatabase()

	token := c.Param("token")
	urlObject, err := repository.FetchUrlDBObject(db, token)

	if err != nil || urlObject == nil {
		logger.Error().Err(err).Msg("Failed to fetch url")
		return response.InvalidUrlsProvided.ReturnResponse(c, nil)
	}

	originalUrl := urlObject.URL
	return c.Redirect(http.StatusTemporaryRedirect, originalUrl)
}
