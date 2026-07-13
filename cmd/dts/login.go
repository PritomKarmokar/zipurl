package dts

import "time"

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=5,max=72"`
}

type UserLoginResponse struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	TokenType    string     `json:"token_type"`
	LastLoginAt  *time.Time `json:"last_login_at"`
}
