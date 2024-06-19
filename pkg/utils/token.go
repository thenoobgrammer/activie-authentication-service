package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserID        string `json:"userId" validate:"required"`
	Email         string `json:"email" validate:"required"`
	EmailVerified bool   `json:"emailVerified" validate:"required"`
	Scopes        bool   `json:"scopes,omitempty" validate:"required"`
}

func GetClaims(tokenString string, secretKey string) (*UserClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)

	return &UserClaims{
		UserID:        claims["user_id"].(string),
		Email:         claims["email"].(string),
		EmailVerified: claims["email_verified"].(bool),
	}, nil
}
func GenerateToken(claims UserClaims, secretKey string) (*string, error) {
	if claims.Email == "" || claims.UserID == "" {
		return nil, errors.New("invalid claims")
	}
	token := jwt.New(jwt.SigningMethodHS256)

	mapClaims := token.Claims.(jwt.MapClaims)

	mapClaims["user_id"] = claims.UserID
	mapClaims["email"] = claims.Email
	mapClaims["email_verified"] = claims.EmailVerified
	mapClaims["exp"] = time.Now().AddDate(1, 0, 0).Unix()

	tokenString, err := token.SignedString(secretKey)

	return &tokenString, err
}
