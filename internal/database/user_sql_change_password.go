package database

import (
	"auth-service/pkg/constants"
	"auth-service/pkg/utils"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func ChangePassword(email string, userID string, uhCurr string, uhNew string) error {
	uint64UserID, err := utils.StringToUint64(userID)
	if err != nil {
		utils.LogError("UpdatePassword", constants.ERROR_DURING_ROW_SCAN, err)
		return err
	}
	// retrieve old hashed pwd based on user infos
	var hCurr string
	if err = GetClient().QueryRow("SELECT password FROM users WHERE id = ? AND email = ?", uint64UserID, email).Scan(&hCurr); err != nil {
		utils.LogError("UpdatePassword", constants.ERROR_DURING_ROW_SCAN, err)
		return err
	}

	// we compare both hash and unhashed current pwds
	if err = bcrypt.CompareHashAndPassword([]byte(hCurr), []byte(uhCurr)); err != nil {
		utils.LogError("UpdatePassword", constants.ERROR_DURING_ROW_SCAN, err)
		return err
	}

	// on success of previous method, we generate a new hashed pwd
	hNew, err := bcrypt.GenerateFromPassword([]byte(uhNew), HASH_SALT)
	if err != nil {
		utils.LogError("UpdatePassword", constants.ERROR_DURING_ROW_SCAN, err)
		return err
	}

	// update current pwd
	rows, err := GetClient().Exec("UPDATE users SET password = ? WHERE id = ? AND email = ?", hNew, uint64UserID, email)
	if err != nil {
		log.Println("RAN UPDATED")
		utils.LogError("UpdatePassword", constants.ERROR_DURING_ROW_SCAN, err)
		return err
	}

	// check if there was any affected row within the db
	affectedRows, err := rows.RowsAffected()
	if err != nil {
		utils.LogError("UpdatePassword", constants.ERROR_DURING_ROW_SCAN, err)
		return err
	}

	if affectedRows <= 0 {
		utils.LogError("UpdatePassword", constants.ERROR_DURING_ROW_SCAN, err)
		return errors.New("current password do not match")
	}

	return nil
}
