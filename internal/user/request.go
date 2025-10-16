package user

type NewUserRequest struct {
	AccountType *string
	Avatar      *string
	Email       string
	GoogleID    *string
	FullName    string
	Password    *string
	SystemRoles *string
}

type GoogleSignupRequest struct {
	Avatar   string `json:"avatar" validate:"required"`
	Email    string `json:"email" validate:"required"`
	GoogleID string `json:"googleId" validate:"required"`
	FullName string `json:"fullName" validate:"required"`
}

type SignupRequest struct {
	Email         string  `json:"email" validate:"required"`
	GoogleID      *string `json:"googleId" validate:"required"`
	FullName      *string `json:"fullName" validate:"required"`
	Password      string  `json:"password" validate:"required"`
	PreferredCity *string `json:"preferredCity" validate:"required"`
}
