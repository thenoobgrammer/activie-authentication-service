package database_session

import (
	"auth-service/internal/database"
	"auth-service/pkg/constants"
	"auth-service/pkg/utils"
)

func DeleteFromUserId(userID string) bool {
	const FUNC_NAME = "DeleteFromUserId"

	uint64UserID, err := utils.StringToUint64(userID)
	if err != nil {
		utils.LogError(FUNC_NAME, constants.ERROR_DURING_CONVERSION, err)
		return false
	}

	query := `DELETE FROM user_sessions WHERE user_id = $1`

	result, err := database.GetClient().Exec(query, uint64UserID)
	if err != nil {
		utils.LogError(FUNC_NAME, constants.ERROR_DURING_DELETE, err)
		return false
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected <= 0 {
		utils.LogWarn(FUNC_NAME, constants.WARNING_NO_ROWS_AFFECTED, nil)
		return false
	}

	return true
}

func DeleteFromToken(token string) bool {
	const FUNC_NAME = "DeleteFromToken"

	query := `DELETE FROM user_sessions WHERE token = $1`

	result, err := database.GetClient().Exec(query, token)
	if err != nil {
		utils.LogError(FUNC_NAME, constants.ERROR_DURING_DELETE, err)
		return false
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected <= 0 {
		utils.LogWarn(FUNC_NAME, constants.WARNING_NO_ROWS_AFFECTED, nil)
		return false
	}
	return true
}

func DeleteFromId(sessionId string) bool {
	const FUNC_NAME = "DeleteFromId"

	query := `DELETE FROM user_sessions WHERE id = $1`

	result, err := database.GetClient().Exec(query, sessionId)
	if err != nil {
		utils.LogError(FUNC_NAME, constants.ERROR_DURING_DELETE, err)
		return false
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected <= 0 {
		utils.LogWarn(FUNC_NAME, constants.WARNING_NO_ROWS_AFFECTED, nil)
		return false
	}
	return true
}
