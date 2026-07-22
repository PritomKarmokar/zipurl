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

	const anonymousClickCount int64 = 1000000

	authHeader := c.Request().Header.Get("Authorization")
	var user *model.User

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

	urlObject, err := repository.FindExistingURL(db, reqBody.Url)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch url")
		return response.TechnicalError.ReturnResponse(c, nil)
	}
	if urlObject != nil {
		logger.Info().Msgf("short url already exists for this url: %s", reqBody.Url)
		shortUrl := viper.GetString("ZIP_URL_BASE_URL") + "/" + urlObject.HashedToken
		responseData := map[string]interface{}{
			"short_url": shortUrl,
		}
		return response.GenericSuccess200.ReturnResponse(c, responseData)
	}
	// Build the short URL
	tokenId := utils.GenerateULID()
	currentTime := time.Now()
	uniqueToken := utils.EncodeString(tokenId[20:]) // Base62 encoded token from last 6 chars of ULID

	newUrl := &model.URL{
		ID:          tokenId,
		URL:         reqBody.Url,
		HashedToken: uniqueToken,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
	}

	if user != nil {
		newUrl.UserID = &user.ID
	} else {
		newUrl.ClickCount = anonymousClickCount
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
