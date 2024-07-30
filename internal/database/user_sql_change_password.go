package database

import (
	"auth-service/pkg/constants"
	"auth-service/pkg/utils"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func ChangePassword(email string, userID string, uhCurr string, uhNew string) (bool, error) {
	uint64UserID, err := utils.StringToUint64(userID)
	if err != nil {
		utils.LogError("ChangePassword", constants.ERROR_DURING_ROW_SCAN, err)
		return false, errors.New("invalid user ID")
	}
	// retrieve old hashed pwd based on user infos
	var hCurr string
	if err = GetClient().QueryRow("SELECT password FROM users WHERE id = ? AND email = ?", uint64UserID, email).Scan(&hCurr); err != nil {
		if err == sql.ErrNoRows {
			utils.LogError("ChangePassword", "User not found or incorrect email", err)
			return false, errors.New("user not found or incorrect email")
		}
		utils.LogError("ChangePassword", constants.ERROR_DURING_ROW_SCAN, err)
		return false, errors.New("error retrieving current password")
	}

	// we compare both hash and unhashed current pwds
	if err = bcrypt.CompareHashAndPassword([]byte(hCurr), []byte(uhCurr)); err != nil {
		utils.LogError("ChangePassword", constants.ERROR_DURING_ROW_SCAN, err)
		return false, nil
	}

	// on success of previous method, we generate a new hashed pwd
	hNew, err := bcrypt.GenerateFromPassword([]byte(uhNew), HASH_SALT)
	if err != nil {
		utils.LogError("ChangePassword", constants.ERROR_DURING_ROW_SCAN, err)
		return false, errors.New("error generating new password")
	}

	// update current pwd
	rows, err := GetClient().Exec("UPDATE users SET password = ? WHERE id = ? AND email = ?", string(hNew), uint64UserID, email)
	if err != nil {
		utils.LogError("ChangePassword", constants.ERROR_DURING_ROW_SCAN, err)
		return false, errors.New("error updating password")
	}

	// check if there was any affected row within the db
	affectedRows, err := rows.RowsAffected()
	if err != nil {
		utils.LogError("ChangePassword", constants.ERROR_DURING_ROW_SCAN, err)
		return false, errors.New("error checking update status")
	}

	if affectedRows <= 0 {
		utils.LogError("ChangePassword", constants.ERROR_DURING_ROW_SCAN, err)
		return false, errors.New("password not changed, please try again")
	}

	return true, nil
}
