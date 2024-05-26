package service

import (
	"auth-service/pkg/models"
	"auth-service/pkg/utils"
)

func RequestToken(claims models.UserClaims) (*string, error) {
	return utils.GenerateToken(claims)
}
func GetClaims(token string) (*models.UserClaims, error) {
	return utils.GetClaims(token)
}
func RefreshToken(token string) (*string, error) {
	userClaims, err := utils.GetClaims(token)
	if err != nil {
		return nil, err
	}

	newToken, err := utils.GenerateToken(*userClaims)

	return newToken, err
}
