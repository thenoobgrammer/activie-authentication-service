package token

import "time"

type RefreshToken struct {
	ID        string    `json:"id"`
	TokenHash string    `json:"tokenHash"`
	UserID    *string   `json:"userId"`
	ClientID  *string   `json:"clientId"`
	ExpiresAt time.Time `json:"expiresAt"`
	IsRevoked bool      `json:"isRevoked"`
	CreatedAt time.Time `json:"createdAt"`
}
type Token struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Principal    string `json:"principal"`
}

type AnonymousToken struct {
	SessionID string `json:"session_id"`
	ExpiresIn int    `json:"expires_in"`
}
