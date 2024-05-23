package service

import (
	sqlutils "auth-service/internal/sql"
	"auth-service/pkg/models"
	"auth-service/pkg/utils"
	"fmt"

	"github.com/go-playground/validator"
)

func validateStruct(req models.CreateUserRequest) error {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' validation failed on tag '%s'", err.Field(), err.Tag()))
		}
		return fmt.Errorf("validation errors:\n%v", errorMessages)
	}
	return nil
}
func ValidateFromBearer(bearer string) (bool, error) {
	userClaims, err := utils.GetClaims(bearer)
	if err != nil {
		return false, err
	}
	return sqlutils.GetByClaims(userClaims)
}
func ValidateCredentials(req models.SystemUserCrendetialsRequest) (bool, error) {
	return sqlutils.GetByCredentials(req.Email, req.Password)
}
