package handler

import (
	"github.com/PritomKarmokar/zipurl/cmd/config"
	"github.com/PritomKarmokar/zipurl/cmd/dts"
	"github.com/PritomKarmokar/zipurl/cmd/response"
	"github.com/labstack/echo/v5"
)

func UrlShortenerHandler(c *echo.Context) error {
	logger := config.GetRequestLogger(c)
	db := config.GetDatabase()
	_ = db

	reqBody := dts.ShortUrlRequest{}
	if err := c.Bind(&reqBody); err != nil {
		logger.Error().Err(err).Msg("failed to bind request body")
		return response.TechnicalError400.ReturnResponse(c, nil)
	}

	if err := c.Validate(reqBody); err != nil {
		logger.Error().Err(err).Msg("Invalid request body")
		return response.DataValidationErr400.ReturnResponse(c, nil)
	}

	logger.Info().Msgf("Shortening URL %v", reqBody.Url)

	return response.GenericSuccess200.ReturnResponse(c, nil)
}
