package models

type Session struct {
	ID          string  `json:"id"`
	AccountType string  `json:"accountType"`
	DeviceType  *string `json:"deviceType"`
	Expiration  string  `json:"exp"`
	LastIP      string  `json:"lastIp"`
	StartTime   string  `json:"start"`
	Token       string  `json:"token"`
	UserEmail   string  `json:"userEmail"`
	UserID      string  `json:"userId"`
}

type SessionRequirements struct {
	AccountType string `json:"accountType"`
	IP          string `json:"ip"`
	UserEmail   string `json:"userEmail"`
	UserId      string `json:"userId"`
}
