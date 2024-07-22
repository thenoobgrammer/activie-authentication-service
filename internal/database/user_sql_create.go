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

func CreateUser(displayName string, email string, externalID *string, fullName string, lat string, lng string, unhashedPwd string, region string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(unhashedPwd), HASH_SALT)
	if err != nil {
		utils.LogError("CreateUser", constants.ERROR_DURING_ROW_SCAN, err)
		panic(err)
	}

	var accountType string

	if externalID != nil {
		accountType = "external"
	} else {
		accountType = "system"
	}

	insertClause := `
		INSERT INTO users (
			account_type,
			agreed_to_terms,
			display_name,
			email,
			email_verified,
			external_id,
			full_name,
			password,
			permissions,
			preferred_locale,
			preferred_location,
			preferred_region,
			preferred_theme,
			role
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	if _, err = GetClient().Exec(insertClause,
		accountType,
		true,
		displayName,
		email,
		false,
		externalID,
		fullName,
		string(hashedPassword),
		strings.Join(constants.DEFAULT_PERMISSIONS[:], ","),
		"en",
		strings.Join([]string{lat, lng}, ","),
		region,
		"light",
		constants.USER,
	); err != nil {
		utils.LogError("CreateUser", constants.ERROR_DURING_ROW_SCAN, err)
		return err
	}

	return nil
}
