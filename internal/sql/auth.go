package sql

import (
	"auth-service/internal/database"
	"auth-service/pkg/constants"
	"auth-service/pkg/models"
	"auth-service/pkg/utils"
	"errors"
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
func GetByID(userID string) (*models.User, error) {
	uint64UserID, err := utils.StringToUint64(userID)
	if err != nil {
		return nil, errors.New("GetByID - Error during conversion (Malformed ID)")
	}

	var user models.User

	selectClause := `
		SELECT
			id, account_type, agreed_to_terms, allow_location_tracking, avatar, bio, city, email,
			email_verified, external_id, full_name, favorites, is_inactive, inactive_date, join_date,
			locale_region, match_organized_count, match_played_count, permissions,
			phone, preferred_locale, preferred_region, preferred_sport, preferred_theme,
			reliability, role, sexe, show_age, show_email, show_groups, show_phone,
			timezone
		FROM users
		WHERE id = ?
	`

	err = database.GetClient().QueryRow(selectClause, uint64UserID).Scan(
		&user.ID, &user.AccountType, &user.AgreedToTerms, &user.AllowLocationTracking, &user.Avatar, &user.Bio, &user.City, &user.DisplayName, &user.Email,
		&user.EmailVerified, &user.ExternalID, &user.FullName, &user.Favorites, &user.IsInactive, &user.InactiveDate, &user.JoinDate,
		&user.LocaleRegion, &user.MatchOrganizedCount, &user.MatchPlayedCount, &user.Permissions,
		&user.Phone, &user.PreferredLocale, &user.PreferredRegion, &user.PreferredSport, &user.PreferredTheme,
		&user.Reliability, &user.Role, &user.Sexe, &user.ShowAge, &user.ShowEmail, &user.ShowGroups, &user.ShowPhone, &user.Timezone,
	)

	return &user, err
}
func GetByClaims(userClaims *models.UserClaims) (*models.User, error) {
	var user models.User

	selectClause := `
		SELECT
			id, account_type, agreed_to_terms, allow_location_tracking, avatar, bio, city, email,
			email_verified, external_id, full_name, favorites, is_inactive, inactive_date, join_date,
			locale_region, match_organized_count, match_played_count, permissions,
			phone, preferred_locale, preferred_region, preferred_sport, preferred_theme,
			reliability, role, sexe, show_age, show_email, show_groups, show_phone,
			timezone
		FROM users
		WHERE id = ?
		AND email = ?
	`

	err := database.GetClient().QueryRow(selectClause, userClaims.UserID, userClaims.Email).Scan(
		&user.ID, &user.AccountType, &user.AgreedToTerms, &user.AllowLocationTracking, &user.Avatar, &user.Bio, &user.City, &user.Email,
		&user.EmailVerified, &user.ExternalID, &user.FullName, &user.Favorites, &user.IsInactive, &user.InactiveDate, &user.JoinDate,
		&user.LocaleRegion, &user.MatchOrganizedCount, &user.MatchPlayedCount, &user.PermissionsString,
		&user.Phone, &user.PreferredLocale, &user.PreferredRegion, &user.PreferredSport, &user.PreferredTheme,
		&user.Reliability, &user.Role, &user.Sexe, &user.ShowAge, &user.ShowEmail, &user.ShowGroups, &user.ShowPhone, &user.Timezone,
	)

	user.IDString = utils.Uint64ToString(user.ID)
	user.Permissions = strings.Split(user.PermissionsString, ",")

	return &user, err
}
func GetByCredentials(email string, unhashedPassword string) (*models.User, error) {
	var hashedPwd string

	err := database.GetClient().QueryRow(`SELECT password FROM users WHERE email = ?`, email).Scan(&hashedPwd)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(unhashedPassword))
	if err != nil {
		return nil, err
	}

	var user models.User

	selectClause := `
		SELECT
			id, account_type, agreed_to_terms, allow_location_tracking, avatar, bio, city, display_name, email,
			email_verified, external_id, full_name, favorites, is_inactive, inactive_date, join_date,
			locale_region, match_organized_count, match_played_count, permissions,
			phone, preferred_locale, preferred_region, preferred_sport, preferred_theme,
			reliability, role, sexe, show_age, show_email, show_groups, show_phone, timezone
		FROM users
		WHERE email = ?
	`

	err = database.GetClient().QueryRow(selectClause, email).Scan(
		&user.ID, &user.AccountType, &user.AgreedToTerms, &user.AllowLocationTracking, &user.Avatar, &user.Bio, &user.City, &user.DisplayName, &user.Email,
		&user.EmailVerified, &user.ExternalID, &user.FullName, &user.Favorites, &user.IsInactive, &user.InactiveDate, &user.JoinDate,
		&user.LocaleRegion, &user.MatchOrganizedCount, &user.MatchPlayedCount, &user.PermissionsString,
		&user.Phone, &user.PreferredLocale, &user.PreferredRegion, &user.PreferredSport, &user.PreferredTheme,
		&user.Reliability, &user.Role, &user.Sexe, &user.ShowAge, &user.ShowEmail, &user.ShowGroups, &user.ShowPhone, &user.Timezone,
	)

	user.IDString = utils.Uint64ToString(user.ID)
	user.Permissions = strings.Split(user.PermissionsString, ",")

	return &user, err
}

func GetGoogleUser(ext_id string, email string) (*models.User, error) {
	var user models.User

	selectClause := `
		SELECT
			id, account_type, agreed_to_terms, allow_location_tracking, avatar, bio, city, email,
			email_verified, external_id, full_name, favorites, is_inactive, inactive_date, join_date,
			locale_region, match_organized_count, match_played_count permissions,
			phone, preferred_locale, preferred_region, preferred_sport, preferred_theme,
			reliability, role, sexe, show_age, show_email, show_groups, show_phone, timezone
		FROM users
		WHERE email = ?
	`

	err := database.GetClient().QueryRow(selectClause, email).Scan(
		&user.ID, &user.AccountType, &user.AgreedToTerms, &user.AllowLocationTracking, &user.Avatar, &user.Bio, &user.City, &user.Email,
		&user.EmailVerified, &user.ExternalID, &user.FullName, &user.Favorites, &user.IsInactive, &user.InactiveDate, &user.JoinDate,
		&user.LocaleRegion, &user.MatchOrganizedCount, &user.MatchPlayedCount, &user.PermissionsString,
		&user.Phone, &user.PreferredLocale, &user.PreferredRegion, &user.PreferredSport, &user.PreferredTheme,
		&user.Reliability, &user.Role, &user.Sexe, &user.ShowAge, &user.ShowEmail, &user.ShowGroups, &user.ShowPhone, &user.Timezone,
	)

	user.IDString = utils.Uint64ToString(user.ID)
	user.Permissions = strings.Split(user.PermissionsString, ",")

	return &user, err
}
