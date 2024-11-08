package database_session

import (
	"auth-service/internal/database"
	"auth-service/pkg/constants"
	"auth-service/pkg/utils"
)

func DeleteFromUserId(userID string) bool {
	uint64UserID, err := utils.StringToUint64(userID)
	if err != nil {
		utils.LogError("DeleteSession", constants.ERROR_DURING_CONVERSION, err)
		return false
	}

	query := `DELETE FROM user_sessions WHERE user_id = ?`

	result, err := database.GetClient().Exec(query, uint64UserID)
	if err != nil {
		utils.LogError("DeleteSession", constants.ERROR_DURING_DELETE, err)
		return false
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected <= 0 {
		utils.LogWarn("DeleteSession", constants.WARNING_NO_ROWS_AFFECTED, nil)
		return false
	}

	return true
}

func DeleteFromToken(token string) bool {
	query := `DELETE FROM user_sessions WHERE token = ?`

	result, err := database.GetClient().Exec(query, token)
	if err != nil {
		utils.LogError("DeleteSession", constants.ERROR_DURING_DELETE, err)
		return false
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected <= 0 {
		utils.LogWarn("DeleteSession", constants.WARNING_NO_ROWS_AFFECTED, nil)
		return false
	}
	return true
}

func DeleteFromId(sessionId string) bool {
	query := `DELETE FROM user_sessions WHERE id = ?`

	result, err := database.GetClient().Exec(query, sessionId)
	if err != nil {
		utils.LogError("DeleteSession", constants.ERROR_DURING_DELETE, err)
		return false
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected <= 0 {
		utils.LogWarn("DeleteSession", constants.WARNING_NO_ROWS_AFFECTED, nil)
		return false
	}
	return true
}
