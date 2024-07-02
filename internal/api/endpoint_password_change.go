package api

import (
	"auth-service/internal/database"
	"auth-service/pkg/api"
	"auth-service/pkg/constants"
	"net/http"
	"unicode"

	"github.com/gin-gonic/gin"
)

type PasswordChangeRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	Email           string `json:"email" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required"`
	UserID          string `json:"userId" validate:"required"`
}

func (req *PasswordChangeRequest) Validate() map[string]string {
	var errors = make(map[string]string)

	if req.Email == "" {
		errors["email"] = "email is missing"
	}
	if req.UserID == "" {
		errors["userId"] = "userId is missing"
	}
	if req.CurrentPassword == "" || req.NewPassword == "" {
		errors["password"] = "password is missing"
	}
	if req.CurrentPassword == req.NewPassword {
		errors["newPassword"] = "your new password cannot be the same as the old one"
	}
	if !isPwdValid(req.NewPassword) {
		errors["newPassword"] = "password needs to have at least one uppercase, lowercase and special character and a minimum lenght of 8"
	}
	return errors
}

func isPwdValid(pwd string) bool {
	var hasUpper, hasLower, hasDigit bool

	if len(pwd) < constants.PASSWORD_MIN_LENGTH {
		return false
	}

	for _, char := range pwd {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

func ChangePassword(g *gin.Context) {
	var req PasswordChangeRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		api.HandleError(g, api.InvalidJSON())
		return
	}

	if errors := req.Validate(); len(errors) > 0 {
		api.HandleError(g, api.InvalidRequestData(errors))
		return
	}

	if err := database.ChangePassword(req.Email, req.UserID, req.CurrentPassword, req.NewPassword); err != nil {
		api.HandleError(g, api.NewInternalServerError(err))
		return
	}

	api.HandleSuccess(g, http.StatusOK, "password changed")
}
