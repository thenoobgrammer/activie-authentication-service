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
func CreateUser(req models.CreateUserRequest) (*models.User, error) {
	err := validateStruct(req)
	if err != nil {
		return nil, err
	}

	insertedID, err := sqlutils.CreateUser(req)
	if err != nil {
		return nil, err
	}

	return sqlutils.GetByID(*insertedID)
}

// func CreateExternalUser(req models.CreateUserRequest) (*models.User, error) {
// 	insertedID, err := sqlutils.CreateUser(req)
// 	if err != nil {
// 		return nil, err
// 	}

//		return sqlutils.GetByID(*insertedID)
//	}
func QueryUserFromBearer(bearer string) (*models.User, error) {
	userClaims, err := utils.GetClaims(bearer)
	if err != nil {
		return nil, err
	}
	return sqlutils.GetByClaims(userClaims)
}
func QueryUserFromCredentials(req models.SystemUserCrendetialsRequest) (*models.User, error) {
	return sqlutils.GetByCredentials(req.Email, req.Password)
}

// TODO
// func QueryExternalUser(req models.ExternalUserCrendetialsRequest) (*models.User, error) {
// 	provider := req.Provider

// 	switch provider {
// 	case "google":
// 		return sqlutils.GetGoogleUser(req.ExternalId, req.Email)
// 	}

// 	return nil, errors.New("unknown provider")
// }
