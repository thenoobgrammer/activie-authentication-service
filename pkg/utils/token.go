package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	AccountType string   `json:"accountType" validate:"required"`
	Email       string   `json:"email" validate:"required"`
	UserID      string   `json:"userId" validate:"required"`
	UserRoles   []string `json:"userRoles" validate:"required"`
}

func GetClaims(tokenString string, secretKey []byte) *UserClaims {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil
	}

	claims := token.Claims.(jwt.MapClaims)
	return &UserClaims{
		AccountType: claims["accounType"].(string),
		Email:       claims["email"].(string),
		UserID:      claims["userId"].(string),
		UserRoles:   claims["userRoles"].([]string),
	}
}
func GenerateToken(claims UserClaims, secretKey []byte) *string {
	if claims.Email == "" || claims.UserID == "" {
		return nil
	}
	token := jwt.New(jwt.SigningMethodHS256)

	mapClaims := token.Claims.(jwt.MapClaims)

	mapClaims["accountType"] = claims.UserID
	mapClaims["email"] = claims.Email
	mapClaims["exp"] = time.Now().AddDate(1, 0, 0).Unix()
	mapClaims["userId"] = claims.UserID

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		LogError("GenerateToken", "error during signing token", err)
		return nil
	}

	return &tokenString
}
