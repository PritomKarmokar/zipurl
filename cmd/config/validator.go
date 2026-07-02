package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

// EchoValidator wraps a go-playground validator for Echo.
type EchoValidator struct {
	Validator *validator.Validate
}

// Validate implements echo.Validator.
func (v *EchoValidator) Validate(i interface{}) error {
	return v.Validator.Struct(i)
}

// RegisterValidator attaches the validator to the Echo instance.
// Call RegisterValidator(e) from `EchoConfig` (e.g., after middleware setup).
func RegisterValidator(e *echo.Echo) {
	e.Validator = &EchoValidator{Validator: validator.New()}

	// Example: register custom validations if needed
	// _ = e.Validator.(*EchoValidator).Validator.RegisterValidation("myrule", myRuleFunc)
}
