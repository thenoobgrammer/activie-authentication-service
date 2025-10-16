package session

import (
	"auth-service/pkg/env"
	"auth-service/pkg/logs"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	Email     string   `json:"email,omitempty"`
	UserID    string   `json:"userId,omitempty"`
	UserRoles []string `json:"userRoles" validate:"required"`
	ClientID  string   `json:"clientId,omitempty"`
	jwt.RegisteredClaims
}

type ClaimsParams struct {
	UserID    string   `json:"userId"`
	UserRoles []string `json:"userRoles"`
	Email     string   `json:"email"`
	ClientID  string   `json:"clientId"`
}

func NewClaims(params ClaimsParams) Claims {
	jti := uuid.New().String()

	claims := Claims{
		Email:     params.Email,
		UserID:    params.UserID,
		UserRoles: params.UserRoles,
		ClientID:  params.ClientID,
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

func GenerateJWT(claims jwt.Claims) (*string, error) {
	key := []byte(env.TOKEN_SECRET)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(key)
	if err != nil {
		logs.Error("GenerateJWT", "error during signing token", err)
		return nil, err
	}

	return &tokenString, nil
}

func GenerateOpaqueToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func ExtractClaims(tokenStr string) (*Claims, error) {
	var claims Claims
	extractedClaims, err := ExtractJWTClaims(tokenStr, &claims)
	if err != nil {
		return nil, err
	}
	return extractedClaims, nil
}

func ExtractJWTClaims[T jwt.Claims](tokenStr string, claims T) (T, error) {
	key := []byte(env.TOKEN_SECRET)

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected.signing.method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return claims, err
	}

	validatedClaims, ok := token.Claims.(T)
	if !ok || !token.Valid {
		return claims, fmt.Errorf("invalid.token.claims")
	}

	exp, err := validatedClaims.GetExpirationTime()
	if err != nil {
		return claims, fmt.Errorf("failed.to.get.token.expiration")
	}

	if exp.Before(time.Now()) {
		return claims, fmt.Errorf("token.expired")
	}

	return validatedClaims, nil
}
