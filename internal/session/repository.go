package session

import (
	"auth-service/internal/infra/database"
	"auth-service/pkg/constants"
	"auth-service/pkg/logs"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	CreateSession(params NewSessionParams) error
	CreateAnonymousSession(expiresAt time.Time) (*AnonymousSession, error)

	GetSessionByJTI(jti string) (*Session, error)
	GetSessionByUserID(userID string) (*Session, error)
	GetSessionByClientID(clientID string) (*Session, error)

	DeactivateSession(jti string) error

	DeleteSessionByJTI(jti string) error
	DeleteSessionByUserID(userID string) error
	DeleteSessionByClientID(clientID string) error
}

type repo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repo{db: db}
}

/* ==== ========== INSERT ================== =====================*/

func (r *repo) CreateSession(params NewSessionParams) error {
	const FuncName = "CreateSession"

	query := `INSERT INTO authenticated_sessions (session_id, user_id, token_jti, start_time, exp, last_ip, device_type, is_active) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	args := []any{params.SessionID, params.UserID, params.TokenJTI, params.StartTime, params.Exp, params.LastIP, params.DeviceType, params.IsActive}

	result, err := database.GetClient().Exec(query, args...)
	if err != nil {
		logs.Error(FuncName, constants.ERROR_DURING_INSERT, err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		logs.Warn(FuncName, constants.WARNING_NO_ROWS_AFFECTED, nil)
		return nil
	}

	return nil
}

func (r *repo) CreateAnonymousSession(expiresAt time.Time) (*AnonymousSession, error) {
	query := `INSERT INTO anonymous_sessions (session_id, expires_at) VALUES ($1,$2) RETURNING *`
	args := []any{uuid.New().String(), expiresAt}

	var sess AnonymousSession

	if err := database.GetClient().QueryRow(query, args...).Scan(
		&sess.SessionID,
		&sess.CreatedAt,
		&sess.ExpiresAt,
	); err != nil {
		return nil, err
	}

	return &sess, nil
}

/* ==== ========== GET ================== =====================*/

func (r *repo) GetSessionByJTI(jti string) (*Session, error) {
	return querySession("token_jti", []any{jti})
}

func (r *repo) GetSessionByUserID(userId string) (*Session, error) {
	return querySession("user_id", []any{userId})
}

func (r *repo) GetSessionByClientID(clientId string) (*Session, error) {
	return querySession("client_id", []any{clientId})
}

/* ==== ========== Deactivate ================== =====================*/

func (r *repo) DeactivateSession(jti string) error {
	const FuncName = "DeactivateSession"

	queryUser := `UPDATE authenticated_sessions SET is_active = false WHERE token_jti = $1`
	result, err := database.GetClient().Exec(queryUser, jti)
	if err != nil {
		logs.Error(FuncName, constants.ERROR_DURING_DELETE, err)
		return nil
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected <= 0 {
		logs.Warn(FuncName, constants.WARNING_NO_ROWS_AFFECTED, nil)
	}

	return nil
}

/* ==== ========== DELETE ================== =====================*/

func (r *repo) DeleteSessionByUserID(userID string) error {
	return deleteSession("user_id", []any{userID})
}

func (r *repo) DeleteSessionByJTI(jti string) error {
	return deleteSession("token_jti", []any{jti})
}

func (r *repo) DeleteSessionByClientID(clientId string) error {
	return deleteSession("client_id", []any{clientId})
}

func querySession(key string, args []any) (*Session, error) {
	var session Session

	queryUser := fmt.Sprintf(`SELECT * FROM authenticated_sessions WHERE %s = $1`, key)

	err := database.GetClient().QueryRow(queryUser, args).Scan(
		&session.SessionID,
		&session.UserID,
		&session.Jti,
		&session.StartTime,
		&session.Expiration,
		&session.LastIP,
		&session.DeviceType,
		&session.IsActive,
	)
	if err == nil {
		return &session, nil
	}

	return &session, nil
}

func deleteSession(key string, args []any) error {
	const FuncName = "deleteSession"

	queryUser := fmt.Sprintf(`DELETE FROM authenticated_sessions WHERE %s = $1`, key)
	result, err := database.GetClient().Exec(queryUser, args...)
	if err != nil {
		logs.Error(FuncName, constants.ERROR_DURING_DELETE, err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected <= 0 {
		logs.Warn(FuncName, constants.WARNING_NO_ROWS_AFFECTED, nil)
	}

	return nil
}
