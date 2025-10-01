package models

type UserSession struct {
	SessionID  string `json:"sessionId"`
	UserID     string `json:"userId"`
	Jti        string `json:"jti"`
	StartTime  string `json:"start"`
	Expiration string `json:"exp"`
	LastIP     string `json:"lastIp"`
	DeviceType string `json:"deviceType"`
	IsActive   bool   `json:"isActive"`
}

type AnonymousSession struct {
	SessionID  string `json:"sessionId"`
	ClientID   string `json:"clientId"`
	Jti        string `json:"jti"`
	StartTime  string `json:"start"`
	Expiration string `json:"exp"`
	LastIP     string `json:"lastIp"`
	DeviceType string `json:"deviceType"`
	IsActive   bool   `json:"isActive"`
}
