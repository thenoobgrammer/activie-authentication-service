package session

import "time"

type ClaimRequirementsParams struct {
	UserID     string   `json:"userId"`
	UserRoles  []string `json:"userRoles"`
	Email      string   `json:"email"`
	ClientID   string   `json:"clientId"`
	LastIP     string   `json:"lastIp"`
	DeviceType string   `json:"deviceType"`
}

type NewSessionParams struct {
	SessionID  string
	ClientID   string
	UserID     string
	TokenJTI   string
	LastIP     string
	DeviceType string
	StartTime  time.Time
	Exp        time.Time
	IsActive   bool
}

type RefreshTokenParams struct {
	RefreshToken string `json:"refresh_token"`
}
