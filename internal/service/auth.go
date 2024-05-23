package service

import (
	sqlutils "auth-service/internal/sql"
	"auth-service/pkg/models"
	"auth-service/pkg/utils"
)

func ValidateID(userID string) bool {
	return sqlutils.GetByID(userID)
}
func ValidateFromBearer(bearer string) bool {
	if bearer == "" {
		return false
	}
	userClaims, err := utils.GetClaims(bearer)
	if err != nil {
		return false
	}
	return sqlutils.GetByID(userClaims.UserID)
}
func ValidateCredentials(req models.SystemUserCrendetialsRequest) bool {
	if req.Email == "" || req.Password == "" {
		return false
	}
	return sqlutils.GetByCredentials(req.Email, req.Password)
}
