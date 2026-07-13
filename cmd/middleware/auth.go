package middleware

import (
	"strings"

	"github.com/PritomKarmokar/zipurl/cmd/config"
	"github.com/PritomKarmokar/zipurl/cmd/response"
	"github.com/labstack/echo/v5"
	"github.com/spf13/viper"
)

// JwtTokenAuth middleware for all routes
func JwtTokenAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			requestLogger := config.GetRequestLogger(c)
			clients := config.GetClients()

			// Get token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				requestLogger.Info().Msg("No Authorization header found")
				return response.UnAuthorized.ReturnResponse(c, nil)
			}

			// Extract token (format: "Bearer <token>")
			var token string
			authTokenSegments := strings.Split(authHeader, " ")
			if len(authTokenSegments) != 2 {
				requestLogger.Info().Msg("Invalid Authorization header format")
				return response.UnAuthorized.ReturnResponse(c, nil)
			}

			tokenType := authTokenSegments[0]
			if tokenType == viper.GetString("JWT_AUTH_HEADER_TYPE") {
				token = authTokenSegments[1]
			} else {
				requestLogger.Info().Msg("Invalid Authorization header format")
				return response.UnAuthorized.ReturnResponse(c, nil)
			}

			// Validate token
			claims, err := clients.TokenClient.VerifyToken(token, "access")
			if err != nil {
				requestLogger.Error().Err(err).Msg("Invalid or expired token")
				return response.UnAuthorized.ReturnResponse(c, nil)
			}

			// SECURITY FIX: Safely extract data from claims with type checking
			if data, ok := claims["data"].(map[string]interface{}); ok {
				// Store user info in context
				for key, value := range data {
					if value != nil {
						if strValue, ok := value.(string); ok {
							requestLogger.Debug().Str("key", strValue).Msg("Token claims")
							c.Set(key, strValue)
						} else {
							c.Set(key, value)
						}
					}
				}
			} else {
				requestLogger.Warn().Msg("Invalid claims or data structure")
				return response.UnAuthorized.ReturnResponse(c, nil)
			}
			c.Set("claims", claims)

			return next(c)
		}
	}
}
