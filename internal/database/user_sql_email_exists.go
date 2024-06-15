package database

import (
	"auth-service/pkg/constants"
	"auth-service/pkg/utils"
)

func EmailExists(email string) (bool, error) {
	selectClause := "SELECT COUNT(id) FROM users WHERE email = ?"

	var count int

	if err := GetClient().QueryRow(selectClause, email).Scan(&count); err != nil {
		utils.LogError("GetUserByEmail", constants.ERROR_DURING_ROW_SCAN, err)
		return false, err
	}

	return count > 0, nil
}
