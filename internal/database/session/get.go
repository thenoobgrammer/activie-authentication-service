package database_session

import (
	"auth-service/internal/database"
	"auth-service/internal/models"
	"auth-service/pkg/constants"
	"auth-service/pkg/utils"
)

func Retrieve(sessionId string) *models.Session {
	query := `SELECT * FROM user_sessions WHERE id = ?`

	var session models.Session

	if err := database.GetClient().QueryRow(query, sessionId).Scan(
		&session.ID,
		&session.AccountType,
		&session.DeviceType,
		&session.UserEmail,
		&session.Expiration,
		&session.LastIP,
		&session.StartTime,
		&session.Token,
		&session.UserID,
	); err != nil {
		utils.LogError("Retrieve", constants.ERROR_DURING_ROW_SCAN, err)
		return nil
	}

	return &session
}

func RetrieveFromToken(token string) *models.Session {
	query := `SELECT * FROM user_sessions WHERE token = ?`

	var session models.Session

	if err := database.GetClient().QueryRow(query, token).Scan(
		&session.ID,
		&session.AccountType,
		&session.DeviceType,
		&session.UserEmail,
		&session.Expiration,
		&session.LastIP,
		&session.StartTime,
		&session.Token,
		&session.UserID,
	); err != nil {
		utils.LogError("Retrieve", constants.ERROR_DURING_ROW_SCAN, err)
		return nil
	}

	return &session
}
