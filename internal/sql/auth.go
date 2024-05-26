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
func ValidID(userID string) bool {
	uint64UserID, err := utils.StringToUint64(userID)
	if err != nil {
		return false
	}

	var exists bool

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE id = ?
		)
	`

	err = database.GetClient().QueryRow(query, uint64UserID).Scan(&exists)

	return exists && err == nil
}
func ValidCredentials(email string, unhashedPassword string) (bool, bool, *string) {
	var userID uint64
	var emailVerified bool
	var hashedPwd string

	err := database.GetClient().QueryRow(`SELECT id, email_verified, password FROM users WHERE email = ?`, email).Scan(&userID, &emailVerified, &hashedPwd)
	if err != nil {
		return false, false, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(unhashedPassword))

	userIDString := utils.Uint64ToString(userID)

	return err == nil, emailVerified, &userIDString
}
