package utils

import (
	"auth-service/internal/vault"
	"auth-service/pkg/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(vault.Envars["TOKEN_SECRET"].(string))

func GetClaims(tokenString string) (*models.UserClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)

	return &models.UserClaims{
		UserID:        claims["user_id"].(string),
		Email:         claims["email"].(string),
		EmailVerified: claims["email_verified"].(bool),
	}, nil
}
func GenerateToken(claims models.UserClaims) (*string, error) {
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
