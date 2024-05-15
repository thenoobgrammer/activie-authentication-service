package models

import "time"

type User struct {
	ID                    uint64     `json:"-"`
	IDString              string     `json:"id"`
	AccountType           string     `json:"accountType,omitempty"`
	AgreedToTerms         bool       `json:"agreedToTerms,omitempty"`
	AllowLocationTracking bool       `json:"allowLocationTracking"`
	Avatar                *string    `json:"avatar,omitempty"`
	Bio                   *string    `json:"bio,omitempty"`
	City                  *string    `json:"city,omitempty"`
	DisplayName           string     `json:"displayName,omitempty"`
	Email                 string     `json:"email,omitempty"`
	EmailVerified         bool       `json:"emailVerified,omitempty"`
	ExternalID            *string    `json:"externalId,omitempty"`
	Favorites             *string    `json:"favorites,omitempty"`
	FullName              string     `json:"fullName,omitempty"`
	InactiveDate          *time.Time `json:"inactiveDate,omitempty"`
	IsInactive            bool       `json:"isInactive,omitempty"`
	JoinDate              time.Time  `json:"joinDate,omitempty"`
	LocaleRegion          *string    `json:"localeRegion,omitempty"`
	MatchOrganizedCount   int        `json:"matchOrganizedCount,omitempty"`
	MatchPlayedCount      int        `json:"matchPlayedCount,omitempty"`
	Password              string     `json:"-"`
	PermissionsString     string     `json:"-"`
	Permissions           []string   `json:"permissions,omitempty"`
	Phone                 string     `json:"phone,omitempty"`
	PreferredLocale       string     `json:"preferredLocale,omitempty"`
	PreferredRegion       string     `json:"preferredRegion,omitempty"`
	PreferredSport        string     `json:"preferredSport,omitempty"`
	PreferredTheme        string     `json:"preferredTheme,omitempty"`
	Reliability           int        `json:"reliability,omitempty"`
	Role                  string     `json:"role,omitempty"`
	Sexe                  string     `json:"sexe,omitempty"`
	ShowAge               bool       `json:"showAge,omitempty"`
	ShowEmail             bool       `json:"showEmail,omitempty"`
	ShowGroups            bool       `json:"showGroups,omitempty"`
	ShowPhone             bool       `json:"showPhone,omitempty"`
	Timezone              string     `json:"timezone,omitempty"`
}

type SystemUserCrendetialsRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// type ExternalUserCrendetialsRequest struct {
// 	AccessToken string `json:"accessToken" validate:"required"`
// 	Provider    string `json:"provider" validate:"required"`
// }

type CreateUserRequest struct {
	AccountType   string
	Avatar        string `json:"picture,omitempty"`
	AgreedToTerms bool   `json:"agreedToTerms" validate:"required"`
	DisplayName   string `json:"displayName" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
	EmailVerified bool   `json:"verified_email,omitempty"`
	ExternalID    string `json:"id,omitempty"`
	FullName      string `json:"fullName" validate:"required"`
	Password      string `json:"password,omitempty"`
	Phone         string `json:"phone"`
}
