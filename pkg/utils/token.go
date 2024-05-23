package utils

import (
	"auth-service/internal/vault"
	"auth-service/pkg/models"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(vault.Envars["TOKEN_SECRET"].(string))

func IssueToken(claims models.UserClaims) (*string, error) {
	log.Println("claims", claims)
	if claims.Email == "" || claims.UserID == "" {
		return nil, errors.New("claims invalid")
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
func GenerateJWT(userID uint64, username string, email string, emailVerified bool) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["email"] = email
	claims["email_verified"] = emailVerified
	claims["exp"] = time.Now().AddDate(1, 0, 0).Unix()

	tokenString, err := token.SignedString(secretKey)

	return tokenString, err
}
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
