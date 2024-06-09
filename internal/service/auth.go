package service

import (
	"auth-service/internal/models"
	sqlutils "auth-service/internal/sql"
	"auth-service/pkg/utils"
)

func ValidateID(userID string) bool {
	return sqlutils.ValidID(userID)
}
func ValidateFromBearer(bearer string) (bool, *models.UserClaims) {
	userClaims, err := utils.GetClaims(bearer)
	if bearer == "" || err != nil {
		return false, nil
	}
	return true, userClaims
}
func ValidateCredentials(req models.SystemUserCrendetialsRequest) (bool, bool, *string) {
	if req.Email == "" || req.Password == "" {
		return false, false, nil
	}
	return sqlutils.ValidCredentials(req.Email, req.Password)
}
