package session

import (
	"auth-service/internal/infra/database"
	"auth-service/internal/models"
	"auth-service/pkg/constants"
	"auth-service/pkg/logs"
	"database/sql"

	"github.com/google/uuid"
)

type Repository interface {
	CreateUserSession(params SesssionParams) *string
	CreateAnonymousSession(params AnonymousSessionParams) *string
	GetUserSessionFromId(sessionId string) *models.UserSession
	GetUserSessionFromJti(jti string) *models.UserSession
	GetAnonymousSessionFromId(sessionId string) *models.AnonymousSession
	GetAnonymousSessionFromJti(jti string) *models.AnonymousSession
	DeleteUserSessionFromId(sessionId string) bool
	DeleteUserSessionFromUserId(userId string) bool
	DeleteUserSessionFromJti(jti string) bool
	DeleteAnonymousSessionFromId(sessionId string) bool
	DeleteAnonymousSessionFromClientId(clientId string) bool
	DeleteAnonymousSessionFromJti(jit string) bool
}

type repo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repo{db: db}
}

/* ==== ========== INSERT ================== =====================*/

func (r *repo) CreateUserSession(params SesssionParams) *string {
	const FuncName = "CreateUserSession"

	sessionID := uuid.New().String()

	query := `INSERT INTO authenticated_sessions (session_id, user_id, token_jti, start_time, exp, last_ip, device_type, is_active) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	result, err := database.GetClient().Exec(query,
		sessionID,
		params.UserID,
		params.TokenJTI,
		params.StartTime,
		params.Exp,
		params.LastIP,
		params.DeviceType,
		params.IsActive,
	)
	if err != nil {
		logs.Error(FuncName, constants.ERROR_DURING_INSERT, err)
		return nil
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		logs.Warn(FuncName, constants.WARNING_NO_ROWS_AFFECTED, nil)
		return nil
	}

	return &sessionID
}

func (r *repo) CreateAnonymousSession(params AnonymousSessionParams) *string {
	const FuncName = "CreateAnonymousSession"

	sessionID := uuid.New().String()

	query := `INSERT INTO anonymous_sessions (session_id, client_id, token_jti, start_time, exp, last_ip, is_active) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	result, err := database.GetClient().Exec(query,
		sessionID,
		params.ClientID,
		params.TokenJTI,
		params.StartTime,
		params.Exp,
		params.LastIP,
		params.IsActive,
	)
	if err != nil {
		logs.Error(FuncName, constants.ERROR_DURING_INSERT, err)
		return nil
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		logs.Warn(FuncName, constants.WARNING_NO_ROWS_AFFECTED, nil)
		return nil
	}

	return &sessionID
}

/* ==== ========== GET ================== =====================*/

func (r *repo) GetUserSessionFromId(sessionId string) *models.UserSession {
	const FuncName = "GetSessionFromId"

	query := `SELECT * FROM authenticated_sessions WHERE session_id = $1`

	var session models.UserSession

	if err := database.GetClient().QueryRow(query, sessionId).Scan(
		&session.SessionID,
		&session.UserID,
		&session.Jti,
		&session.StartTime,
		&session.Expiration,
		&session.LastIP,
		&session.DeviceType,
		&session.IsActive,
	); err != nil {
		logs.Error(FuncName, constants.ERROR_DURING_ROW_SCAN, err)
		return nil
	}

	return &session
}

func (r *repo) GetUserSessionFromJti(jti string) *models.UserSession {
	const FuncName = "GetSessionFromJti"

	query := `SELECT * FROM authenticated_sessions WHERE token_jti = $1`

	var session models.UserSession

	if err := database.GetClient().QueryRow(query, jti).Scan(
		&session.SessionID,
		&session.UserID,
		&session.Jti,
		&session.StartTime,
		&session.Expiration,
		&session.LastIP,
		&session.DeviceType,
		&session.IsActive,
	); err != nil {
		logs.Error(FuncName, constants.ERROR_DURING_ROW_SCAN, err)
		return nil
	}

	return &session
}

func (r *repo) GetAnonymousSessionFromId(sessionId string) *models.AnonymousSession {
	const FuncName = "GetAnonymousSessionFromId"

	query := `SELECT * FROM anonymous_sessions WHERE session_id = $1`

	var session models.AnonymousSession

	if err := database.GetClient().QueryRow(query, sessionId).Scan(
		&session.SessionID,
		&session.ClientID,
		&session.Jti,
		&session.StartTime,
		&session.Expiration,
		&session.LastIP,
		&session.DeviceType,
		&session.IsActive,
	); err != nil {
		logs.Error(FuncName, constants.ERROR_DURING_ROW_SCAN, err)
		return nil
	}

	return &session
}

func (r *repo) GetAnonymousSessionFromJti(jti string) *models.AnonymousSession {
	const FuncName = "GetAnonymousSessionFromId"

	query := `SELECT * FROM anonymous_sessions WHERE token_jti = $1`

	var session models.AnonymousSession

	if err := database.GetClient().QueryRow(query, jti).Scan(
		&session.SessionID,
		&session.ClientID,
		&session.Jti,
		&session.StartTime,
		&session.Expiration,
		&session.LastIP,
		&session.DeviceType,
		&session.IsActive,
	); err != nil {
		logs.Error(FuncName, constants.ERROR_DURING_ROW_SCAN, err)
		return nil
	}

	return &session
}

/* ==== ========== DELETE ================== =====================*/

func (r *repo) DeleteUserSessionFromId(sessionId string) bool {
	const FuncName = "DeleteUserSessionFromUserId"

	query := `DELETE FROM authenticated_sessions WHERE session_id = $1`

	result, err := database.GetClient().Exec(query, sessionId)
	if err != nil {
		logs.Error(FuncName, constants.ERROR_DURING_DELETE, err)
		return false
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected <= 0 {
		logs.Warn(FuncName, constants.WARNING_NO_ROWS_AFFECTED, nil)
		return false
	}

	return true
}

func (r *repo) DeleteUserSessionFromUserId(userID string) bool {
	const FuncName = "DeleteUserSessionFromUserId"

	query := `DELETE FROM authenticated_sessions WHERE user_id = $1`

	result, err := database.GetClient().Exec(query, userID)
	if err != nil {
		logs.Error(FuncName, constants.ERROR_DURING_DELETE, err)
		return false
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected <= 0 {
		logs.Warn(FuncName, constants.WARNING_NO_ROWS_AFFECTED, nil)
		return false
	}

	return true
}

func (r *repo) DeleteUserSessionFromJti(jit string) bool {
	const FuncName = "DeleteUserSessionFromJit"

	query := `DELETE FROM authenticated_sessions WHERE token_jti = $1`

	result, err := database.GetClient().Exec(query, jit)
	if err != nil {
		logs.Error(FuncName, constants.ERROR_DURING_DELETE, err)
		return false
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected <= 0 {
		logs.Warn(FuncName, constants.WARNING_NO_ROWS_AFFECTED, nil)
		return false
	}
	return true
}

func (r *repo) DeleteAnonymousSessionFromId(sessionId string) bool {
	const FuncName = "DeleteAnonymousSessionFromId"

	query := `DELETE FROM anonymous_sessions WHERE session_id = $1`

	result, err := database.GetClient().Exec(query, sessionId)
	if err != nil {
		logs.Error(FuncName, constants.ERROR_DURING_DELETE, err)
		return false
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected <= 0 {
		logs.Warn(FuncName, constants.WARNING_NO_ROWS_AFFECTED, nil)
		return false
	}
	return true
}

func (r *repo) DeleteAnonymousSessionFromClientId(clientId string) bool {
	const FuncName = "DeleteAnonymousSessionFromClientId"

	query := `DELETE FROM anonymous_sessions WHERE client_id = $1`

	result, err := database.GetClient().Exec(query, clientId)
	if err != nil {
		logs.Error(FuncName, constants.ERROR_DURING_DELETE, err)
		return false
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected <= 0 {
		logs.Warn(FuncName, constants.WARNING_NO_ROWS_AFFECTED, nil)
		return false
	}
	return true
}

func (r *repo) DeleteAnonymousSessionFromJti(jti string) bool {
	const FuncName = "DeleteAnonymousSessionFromJit"

	query := `DELETE FROM anonymous_sessions WHERE token_jti = $1`

	result, err := database.GetClient().Exec(query, jti)
	if err != nil {
		logs.Error(FuncName, constants.ERROR_DURING_DELETE, err)
		return false
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected <= 0 {
		logs.Warn(FuncName, constants.WARNING_NO_ROWS_AFFECTED, nil)
		return false
	}
	return true
}
