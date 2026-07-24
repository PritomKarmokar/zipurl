package handler

import (
	"net/http"
	"strings"
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
	clients := config.GetClients()
	reqBody := dts.ShortUrlRequest{}

	if err := c.Bind(&reqBody); err != nil {
		logger.Error().Err(err).Msg("failed to bind request body")
		return response.TechnicalError.ReturnResponse(c, nil)
	}
	if err := c.Validate(reqBody); err != nil {
		logger.Error().Err(err).Msg("Invalid request body")
		return response.DataValidationErr.ReturnResponse(c, nil)
	}

	var user *model.User
	authHeader := c.Request().Header.Get("Authorization")

	if authHeader != "" {
		// A header was sent, so it must be valid — no silent fallback to anonymous.
		authTokenSegments := strings.Split(authHeader, " ")
		if len(authTokenSegments) != 2 || authTokenSegments[0] != viper.GetString("JWT_AUTH_HEADER_TYPE") {
			logger.Warn().Msg("malformed authorization header")
			return response.UnAuthorized.ReturnResponse(c, nil)
		}
		token := authTokenSegments[1]

		claims, err := clients.TokenClient.VerifyToken(token, "access")
		if err != nil {
			logger.Warn().Err(err).Msg("failed to verify token")
			return response.UnAuthorized.ReturnResponse(c, nil)
		}
		data, ok := claims["data"].(map[string]interface{})
		if !ok {
			logger.Warn().Msg("token claims missing data field")
			return response.UnAuthorized.ReturnResponse(c, nil)
		}
		userId, ok := data["id"].(string)
		if !ok {
			logger.Warn().Msg("token claims missing id field")
			return response.UnAuthorized.ReturnResponse(c, nil)
		}
		fetchedUser, err := repository.GetUserByID(db, userId)
		if err != nil {
			logger.Error().Err(err).Msg("failed to fetch user")
			return response.TechnicalError.ReturnResponse(c, nil)
		}
		if fetchedUser == nil {
			logger.Warn().Msg("user not found for token")
			return response.UnAuthorized.ReturnResponse(c, nil)
		}
		user = fetchedUser
	}

	if user == nil && (reqBody.Expiry != "" || reqBody.MaximumClicks != nil) {
		logger.Warn().Msg("Anonymous User Attempted to set restricted fields")
		return response.PermissionForbidden.ReturnResponse(c, nil)
	}

	var shortUrl string
	currentTime := time.Now()

	// Request Handling for Anonymous User
	if user == nil {
		exisingUrl, err := repository.FindExistingURLForAnonymousUser(db, reqBody.Url)
		if err != nil {
			logger.Error().Err(err).Msg("failed to find existing url")
			return response.TechnicalError.ReturnResponse(c, nil)
		}
		if exisingUrl != nil {
			logger.Info().Msgf("short url already exists for this url: %s", reqBody.Url)
			shortUrl = viper.GetString("ZIP_URL_BASE_URL") + "/" + exisingUrl.HashedToken
		} else {
			// Anonymous links are always unrestricted/shared by design — do not set ExpiresAt or MaxClicks here.
			urlObject := &model.URL{
				ID:          utils.GenerateULID(),
				URL:         reqBody.Url,
				HashedToken: utils.GenerateShortToken(),
				CreatedAt:   currentTime,
				UpdatedAt:   currentTime,
				ClickCount:  0,
			}
			if err := repository.CreateUrlDBObject(db, urlObject); err != nil {
				logger.Warn().Err(err).Msg("Failed to create url db object")
				return response.ShortURLCreationFailed.ReturnResponse(c, nil)
			}
			shortUrl = viper.GetString("ZIP_URL_BASE_URL") + "/" + urlObject.HashedToken
		}
		responseData := map[string]interface{}{
			"short_url": shortUrl,
		}
		return response.GenericSuccess200.ReturnResponse(c, responseData)
	}

	var expiryTime time.Time
	if reqBody.Expiry != "" {
		duration, err := utils.ParseExpiryDuration(reqBody.Expiry)
		if err != nil {
			logger.Warn().
				Str("expiry", reqBody.Expiry).
				Err(err).
				Msg("invalid expiry duration")

			return response.ExpiryTimeDataValidationError.ReturnResponse(c, nil)
		}
		expiryTime = currentTime.Add(duration)
	}

	existingUrl, err := repository.FindExistingUrlForLoggedInUser(db, reqBody.Url, user.ID, currentTime)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch url")
		return response.TechnicalError.ReturnResponse(c, nil)
	}

	if existingUrl != nil {
		logger.Info().Msgf("short url already exists for this url: %s", reqBody.Url)
		shortUrl = viper.GetString("ZIP_URL_BASE_URL") + "/" + existingUrl.HashedToken
	} else {
		// Build the new short URL for Logged Users perspective
		newUrl := &model.URL{
			ID:          utils.GenerateULID(),
			URL:         reqBody.Url,
			HashedToken: utils.GenerateShortToken(),
			UserID:      &user.ID,
			CreatedAt:   currentTime,
			UpdatedAt:   currentTime,
			ClickCount:  0,
		}
		if reqBody.Expiry != "" {
			newUrl.ExpiresAt = &expiryTime
		}
		if reqBody.MaximumClicks != nil {
			newUrl.MaxClicks = reqBody.MaximumClicks
		}

		if err := repository.CreateUrlDBObject(db, newUrl); err != nil {
			logger.Warn().Err(err).Msg("Failed to create url db object")
			return response.ShortURLCreationFailed.ReturnResponse(c, nil)
		}
		logger.Info().Msg("URL Shortener DB object created successfully")
		shortUrl = viper.GetString("ZIP_URL_BASE_URL") + "/" + newUrl.HashedToken
	}

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
