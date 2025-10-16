package token

type NewRefreshTokenRecordParams struct {
	UserID    *string
	ClientID  *string
	TokenHash string
}
