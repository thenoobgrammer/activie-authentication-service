package service

import (
	"auth-service/pkg/models"
	"auth-service/pkg/utils"
)

func IssueToken(claims models.UserClaims) (*string, error) {
	return utils.IssueToken(claims)
}
func GetClaims(token string) (*models.UserClaims, error) {
	return utils.GetClaims(token)
}
func RefreshToken(token string) (*string, error) {
	userClaims, err := utils.GetClaims(token)
	if err != nil {
		return nil, err
	}

	newToken, err := utils.IssueToken(*userClaims)

	return newToken, err
}
