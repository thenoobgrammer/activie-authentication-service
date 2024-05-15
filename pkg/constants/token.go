package constants

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	jwt.Claims
	ID            uint64 `json:"user_id"`
	Email         string `json:"email"`
	Username      string `json:"username"`
	EmailVerified bool   `json:"email_verified"`
}
