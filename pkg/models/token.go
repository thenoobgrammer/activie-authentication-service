package models

type UserClaims struct {
	UserID        string `json:"userId" validate:"required"`
	Email         string `json:"email" validate:"required"`
	EmailVerified bool   `json:"emailVerified" validate:"required"`
}
