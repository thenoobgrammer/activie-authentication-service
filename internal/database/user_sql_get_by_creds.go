package database

import (
	"auth-service/internal/models"
	"auth-service/pkg/constants"
	"auth-service/pkg/utils"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func GetUserByCreds(email string, unhashedPwd string) (*models.User, error) {
	var user models.User

	selectClause := "SELECT password FROM users WHERE email = ?"

	var hashedPwd string

	if err := GetClient().QueryRow(selectClause, email).Scan(&hashedPwd); err != nil {
		utils.LogError("GetUserByCreds", constants.ERROR_DURING_ROW_SCAN, err)
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(unhashedPwd)); err != nil {
		utils.LogError("GetUserByCreds", constants.ERROR_DURING_ROW_SCAN, err)
		return nil, err
	}

	selectClause = `
		SELECT
			id, account_type, agreed_to_terms, allow_location_tracking, avatar, bio, city, email,
			email_verified, external_id, full_name, is_inactive, inactive_date, join_date,
			locale_region, match_organized_count, match_played_count, permissions,
			phone, preferred_locale, preferred_region, preferred_sport, preferred_theme,
			reliability, role, sexe, show_age, show_email, show_groups, show_phone, timezone
		FROM users
		WHERE email = ?
	`

	if err := GetClient().QueryRow(selectClause, email).Scan(
		&user.ID, &user.AccountType, &user.AgreedToTerms, &user.AllowLocationTracking, &user.Avatar, &user.Bio, &user.City, &user.Email,
		&user.EmailVerified, &user.ExternalID, &user.FullName, &user.IsInactive, &user.InactiveDate, &user.JoinDate,
		&user.LocaleRegion, &user.MatchOrganizedCount, &user.MatchPlayedCount, &user.PermissionsString,
		&user.Phone, &user.PreferredLocale, &user.PreferredRegion, &user.PreferredSport, &user.PreferredTheme,
		&user.Reliability, &user.Role, &user.Sexe, &user.ShowAge, &user.ShowEmail, &user.ShowGroups, &user.ShowPhone, &user.Timezone,
	); err != nil {
		utils.LogError("GetUserByCreds", constants.ERROR_DURING_ROW_SCAN, err)
		return nil, err
	}

	user.IDString = utils.Uint64ToString(user.ID)
	user.Permissions = strings.Split(user.PermissionsString, ",")

	return &user, nil
}
