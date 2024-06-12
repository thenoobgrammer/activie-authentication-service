package database

import (
	"auth-service/pkg/constants"
	"auth-service/pkg/utils"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	HASH_SALT = 10
)

func CreateUser(accountType string, avatar *string, displayName string, email string, externalID *string, fullName string, unhashedPwd string, phone *string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(unhashedPwd), HASH_SALT)
	if err != nil {
		utils.LogError("CreateUser", constants.ERROR_DURING_ROW_SCAN, err)
		panic(err)
	}

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
	if _, err = GetClient().Exec(insertClause,
		accountType,
		true,
		avatar,
		displayName,
		email,
		false,
		externalID,
		fullName,
		string(hashedPassword),
		strings.Join(constants.DEFAULT_PERMISSIONS[:], ","),
		phone,
		"en",
		"light",
		constants.USER,
	); err != nil {
		utils.LogError("CreateUser", constants.ERROR_DURING_ROW_SCAN, err)
		return err
	}

	return nil
}
