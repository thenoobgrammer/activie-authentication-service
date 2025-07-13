package database_session

import (
	"auth-service/internal/database"
	"auth-service/internal/models"
	"auth-service/pkg/constants"
	"auth-service/pkg/utils"
	"time"

	"github.com/google/uuid"
)

var (
	SessionExpirationTime = time.Now().AddDate(1, 0, 0)
)

func Create(req models.SessionRequirements, token string) *string {
	const FUNC_NAME = "CreateSession"

	sessionID := uuid.New().String()

	query := `INSERT INTO user_sessions (id, account_type, email, exp, last_ip, token, user_id, user_roles) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	result, err := database.GetClient().Exec(query,
		sessionID,
		req.AccountType,
		req.UserEmail,
		SessionExpirationTime,
		req.IP,
		token,
		req.UserId,
		req.UserRoles,
	)
	if err != nil {
		utils.LogError(FUNC_NAME, constants.ERROR_DURING_INSERT, err)
		return nil
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.LogWarn(FUNC_NAME, constants.WARNING_NO_ROWS_AFFECTED, nil)
		return nil
	}

	return &sessionID
}
