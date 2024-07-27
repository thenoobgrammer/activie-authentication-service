package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	AccountType string `json:"accountType" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Permissions string `json:"permissions" validate:"required"`
	Role        string `json:"role" validate:"required"`
	UserID      string `json:"userId" validate:"required"`
}

func GetClaims(tokenString string, secretKey []byte) (*UserClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	return &UserClaims{
		AccountType: claims["account_type"].(string),
		Email:       claims["email"].(string),
		Permissions: claims["permissions"].(string),
		Role:        claims["role"].(string),
		UserID:      claims["user_id"].(string),
	}, nil
}
func GenerateToken(claims UserClaims, secretKey []byte) (*string, error) {
	if claims.Email == "" || claims.UserID == "" {
		return nil, errors.New("invalid claims")
	}
	token := jwt.New(jwt.SigningMethodHS256)

	mapClaims := token.Claims.(jwt.MapClaims)

	mapClaims["account_type"] = claims.UserID
	mapClaims["email"] = claims.Email
	mapClaims["permissions"] = claims.Permissions
	mapClaims["role"] = claims.Role
	mapClaims["user_id"] = claims.UserID
	mapClaims["exp"] = time.Now().AddDate(1, 0, 0).Unix()

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		LogError("GenerateToken", "error during signing token", err)
		return nil, err
	}

	return &tokenString, nil
}
