package handler

import (
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
		return response.TechnicalError400.ReturnResponse(c, nil)
	}

	if err := c.Validate(reqBody); err != nil {
		logger.Error().Err(err).Msg("Invalid request body provided for signup request")
		return response.DataValidationErr400.ReturnResponse(c, nil)
	}

	user, err := repository.GetUserWithEmailAddress(db, reqBody.Email)
	if user != nil || err != nil {
		logger.Error().Err(err).Msg("User with this email address not found")
		return response.UserAlreadyExistsWithEmail.ReturnResponse(c, nil)
	}

	newUser := &model.User{
		ID:         utils.GenerateULID(),
		UserName:   reqBody.UserName,
		FirstName:  reqBody.FirstName,
		LastName:   reqBody.LastName,
		Email:      reqBody.Email,
		Role:       string(model.UserRole),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		DateJoined: time.Now(),
	}

	if err := newUser.SetPassword(reqBody.Password); err != nil {
		logger.Error().Err(err).Msg("Error occurred while setting user password")
		return response.TechnicalError400.ReturnResponse(c, nil)
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
		"date_joined": utils.FormatTime(newUser.DateJoined),
	}
	return response.UserSignUpSuccess.ReturnResponse(c, responseData)
}
