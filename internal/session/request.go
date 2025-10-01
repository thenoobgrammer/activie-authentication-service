package session

import "time"

type SessionInfoParams struct {
	UserID     string
	UserRoles  []string
	LastIP     string
	DeviceType string
}

type SesssionParams struct {
	SessionID  string
	UserID     string
	TokenJTI   string
	LastIP     string
	DeviceType string
	StartTime  time.Time
	Exp        time.Time
	IsActive   bool
}

type AnonymousSessionInfoParams struct {
	ClientID   string
	LastIP     string
	DeviceType string
}

type AnonymousSessionParams struct {
	ClientID   string
	SessionID  string
	TokenJTI   string
	LastIP     string
	DeviceType string
	StartTime  time.Time
	Exp        time.Time
	IsActive   bool
}
