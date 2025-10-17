package session

import (
	"time"
)

type Session struct {
	SessionID  string `json:"sessionId"`
	ClientID   string `json:"clientId"`
	UserID     string `json:"userId"`
	Jti        string `json:"jti"`
	StartTime  string `json:"start"`
	Expiration string `json:"exp"`
	LastIP     string `json:"lastIp"`
	DeviceType string `json:"deviceType"`
	IsActive   bool   `json:"isActive"`
}

type AnonymousSession struct {
	SessionID string    `json:"sessionId"`
	CreatedAt string    `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}
