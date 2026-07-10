package handler

import (
	"strings"
	"time"

	"github.com/PritomKarmokar/zipurl/cmd/config"
	"github.com/PritomKarmokar/zipurl/cmd/dts"
	"github.com/PritomKarmokar/zipurl/cmd/model"
	"github.com/PritomKarmokar/zipurl/cmd/repository"
	"github.com/PritomKarmokar/zipurl/cmd/response"
	"github.com/PritomKarmokar/zipurl/cmd/utils"
	"github.com/labstack/echo/v5"
)

func UserSignUpHandler(c *echo.Context) error {
	db := config.GetDatabase()
	logger := config.GetRequestLogger(c)

	reqBody := dts.UserSignUp{}
	if err := c.Bind(&reqBody); err != nil {
		logger.Error().Err(err).Msg("failed to bind request body for signup request")
		return response.TechnicalError.ReturnResponse(c, nil)
	}

	reqBody.UserName = strings.TrimSpace(reqBody.UserName)
	reqBody.FirstName = strings.TrimSpace(reqBody.FirstName)
	reqBody.LastName = strings.TrimSpace(reqBody.LastName)
	reqBody.Email = utils.NormalizeEmail(reqBody.Email)

	if err := c.Validate(reqBody); err != nil {
		logger.Error().Err(err).Msg("Invalid request body provided for signup request")
		return response.DataValidationErr.ReturnResponse(c, nil)
	}

	email := reqBody.Email
	userExists, err := repository.UserExistsByEmail(db, email)
	if err != nil {
		logger.Error().Err(err).Msgf("Error Occurred While Fetching User with email %v", email)
		return response.TechnicalError.ReturnResponse(c, nil)
	}
	if userExists {
		logger.Debug().Msgf("User with email %v already exists", email)
		return response.UserAlreadyExistsWithEmail.ReturnResponse(c, nil)
	}

	now := time.Now()

	newUser := &model.User{
		ID:         utils.GenerateULID(),
		UserName:   reqBody.UserName,
		FirstName:  reqBody.FirstName,
		LastName:   reqBody.LastName,
		Email:      email,
		Role:       model.UserRole,
		Status:     model.StatusActive,
		CreatedAt:  now,
		UpdatedAt:  now,
		DateJoined: now,
	}

	if err := newUser.SetPassword(reqBody.Password); err != nil {
		logger.Error().Err(err).Msg("Error occurred while setting user password")
		return response.TechnicalError.ReturnResponse(c, nil)
	}

	if err := repository.CreateNewUserObject(db, newUser); err != nil {
		logger.Error().Err(err).Msg("Error Occurred while creating a new user object")
		return response.UserSignUpFailed.ReturnResponse(c, nil)
	}

	responseData := map[string]any{
		"username":    newUser.UserName,
		"first_name":  newUser.FirstName,
		"last_name":   newUser.LastName,
		"email":       newUser.Email,
		"date_joined": utils.FormatTime(now),
	}
	return response.UserSignUpSuccess.ReturnResponse(c, responseData)
}
