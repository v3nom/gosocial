package facebook

import "time"

// ErrorResponse fb error response
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    int    `json:"code"`
	}
}

// IdentityResponse fb identity response
type IdentityResponse struct {
	ErrorResponse
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// UserResponse fb user info response
type UserResponse struct {
	ErrorResponse
	Email                string `json:"email"`
	Name                 string `json:"name"`
	ID                   string `json:"id"`
	AccessToken          string
	AccessTokenExpiresIn time.Time
}
