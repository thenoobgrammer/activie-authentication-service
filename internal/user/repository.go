package user

import (
	"auth-service/internal/infra/database"
	"database/sql"
)

type Repository interface {
	PeekUserById(userId string) (*User, error)
}

type repo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repo{db: db}
}

func (r *repo) PeekUserById(userID string) (*User, error) {
	var user User
	if err := database.GetClient().QueryRow("SELECT id, email, system_roles FROM users WHERE id=$1", userID).Scan(
		&user.ID,
		&user.Email,
		&user.SystemRoles,
	); err != nil {
		return nil, err
	}
	return &user, nil
}
