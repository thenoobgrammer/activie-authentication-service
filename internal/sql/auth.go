package sql

import (
	"auth-service/internal/database"
	"auth-service/pkg/constants"
	"auth-service/pkg/models"
	"auth-service/pkg/utils"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(req models.CreateUserRequest) (*string, error) {
	insertClause := `
		INSERT INTO users (
			account_type,
			agreed_to_terms,
		 	avatar,
			display_name,
			email,
			email_verified,
			external_id,
			full_name,
			password,
			permissions,
			phone,
			preferred_locale,
			preferred_theme,
			role
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := database.GetClient().Exec(insertClause,
		req.AccountType,
		true,
		req.Avatar,
		req.DisplayName,
		req.Email,
		req.EmailVerified,
		req.ExternalID,
		req.FullName,
		req.Password,
		strings.Join(constants.DEFAULT_PERMISSIONS[:], ","),
		req.Phone,
		"en",
		"light",
		constants.USER,
	)
	if err != nil {
		return nil, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	userID := utils.Int64ToString(lastInsertID)

	return &userID, nil
}
func GetByClaims(userClaims *models.UserClaims) (bool, error) {
	var exists bool

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE id = ?
			AND email = ?
		)
	`

	err := database.GetClient().QueryRow(query, userClaims.UserID, userClaims.Email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
func GetByCredentials(email string, unhashedPassword string) (bool, error) {
	var hashedPwd string

	err := database.GetClient().QueryRow(`SELECT password FROM users WHERE email = ?`, email).Scan(&hashedPwd)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(unhashedPassword))
	if err != nil {
		return false, err
	}

	return true, nil
}
