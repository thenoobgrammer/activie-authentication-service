package session

import (
	"auth-service/pkg/env"
	"auth-service/pkg/logs"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaimInfoParams struct {
	UserID    string   `json:"userId"`
	UserRoles []string `json:"userRoles" validate:"required"`
}

type UserClaims struct {
	UserID    string   `json:"userId"`
	UserRoles []string `json:"userRoles" validate:"required"`
	jwt.RegisteredClaims
}

func NewUserClaims(params UserClaimInfoParams) UserClaims {
	jti := uuid.New().String()

	claims := UserClaims{
		UserID:    params.UserID,
		UserRoles: params.UserRoles,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Subject:   params.UserID,
			Issuer:    "activie.auth.service",
		},
	}

	return claims
}

func GenerateToken(claims UserClaims) *string {
	key := []byte(env.TOKEN_SECRET)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(key)
	if err != nil {
		logs.Error("GenerateToken", "error during signing token", err)
		return nil
	}

	return &tokenString
}

type AnonymousClaimInfoParams struct {
	ClientID string `json:"clientId"`
}

type AnonymousUserClaims struct {
	SessionID string `json:"sessionId"`
	ClientID  string `json:"clientId"`
	IP        string `json:"ip"`
	jwt.RegisteredClaims
}

func NewAnonymousClaims(params AnonymousClaimInfoParams) AnonymousUserClaims {
	sessionID := uuid.New().String()

	jti := uuid.New().String()

	claims := AnonymousUserClaims{
		SessionID: sessionID,
		ClientID:  params.ClientID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Subject:   "anonymous",
			Issuer:    "activie.auth.service",
		},
	}

	return claims
}

func GenerateAnonymousToken(claims AnonymousUserClaims) (*string, error) {
	key := []byte(env.TOKEN_SECRET)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(key)
	if err != nil {
		logs.Error("GenerateAnonymousToken", "error.during.signing.token", err)
		return nil, err
	}

	return &tokenString, nil
}
