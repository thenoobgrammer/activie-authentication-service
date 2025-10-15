package token

import (
	"auth-service/internal/infra/database"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	AddRefreshToken(params NewRefreshTokenRecordParams) error
	RevokeRefreshToken(tokenHash string) (*RefreshToken, error)
	RevokeUserRefreshTokens(userId string) (bool, error)
	GetRefreshToken(tokenHash string) (*RefreshToken, error)
}

type repo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repo{db: db}
}

func (r *repo) AddRefreshToken(params NewRefreshTokenRecordParams) error {
	_, err := database.GetClient().Exec("INSERT INTO refresh_tokens (id, token_hash, user_id, client_id, expires_at) VALUES ($1, $2, $3, $4, $5)",
		uuid.New().String(),
		params.TokenHash,
		*params.UserID,
		*params.ClientID,
		time.Now().AddDate(0, 0, 7),
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) GetRefreshToken(tokenHash string) (*RefreshToken, error) {
	var refreshToken RefreshToken

	if err := database.GetClient().QueryRow("SELECT * from refresh_tokens WHERE token_hash = $1", []any{tokenHash}).Scan(
		&refreshToken.ID,
		&refreshToken.TokenHash,
		&refreshToken.UserID,
		&refreshToken.ClientID,
		&refreshToken.ExpiresAt,
		&refreshToken.IsRevoked,
		&refreshToken.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (r *repo) RevokeRefreshToken(tokenHash string) (*RefreshToken, error) {
	var refreshToken RefreshToken
	if err := database.GetClient().QueryRow("UPDATE refresh_tokens SET is_revoked = true WHERE token_hash = $1 RETURNING *", tokenHash).Scan(
		&refreshToken.ID,
		&refreshToken.TokenHash,
		&refreshToken.UserID,
		&refreshToken.ClientID,
		&refreshToken.ExpiresAt,
		&refreshToken.IsRevoked,
		&refreshToken.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *repo) RevokeUserRefreshTokens(userId string) (bool, error) {
	_, err := database.GetClient().Exec("UPDATE refresh_tokens SET is_revoked = true WHERE user_id = $1", userId)
	if err != nil {
		return false, err
	}
	return true, nil
}
