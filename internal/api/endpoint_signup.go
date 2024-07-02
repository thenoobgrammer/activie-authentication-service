package api

import (
	"auth-service/internal/database"
	"auth-service/internal/vault"
	"auth-service/pkg/api"
	"auth-service/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SignupRequest struct {
	AccountType   string
	Avatar        string `json:"picture,omitempty"`
	AgreedToTerms bool   `json:"agreedToTerms" validate:"required"`
	DisplayName   string `json:"displayName" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
	EmailVerified bool   `json:"verified_email,omitempty"`
	ExternalID    string `json:"id,omitempty"`
	FullName      string `json:"fullName" validate:"required"`
	Password      string `json:"password,omitempty"`
	Phone         string `json:"phone"`
}

func (req *SignupRequest) Validate() map[string]string {
	var errors = make(map[string]string)

	if req.AccountType == "" || (req.AccountType != "external" && req.AccountType != "system") {
		errors["accountType"] = "account type is missing"
	}
	if !req.AgreedToTerms {
		errors["agreedToTerms"] = "user needs to agree to terms"
	}
	if req.Email == "" {
		errors["email"] = "email is required"
	}
	if req.Password == "" {
		errors["password"] = "password is required"
	}

	return errors
}
func Signup(g *gin.Context) {
	var req SignupRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		api.HandleError(g, api.InvalidJSON())
		return
	}

	if errors := req.Validate(); len(errors) > 0 {
		api.HandleError(g, api.InvalidRequestData(errors))
		return
	}

	exists, err := database.EmailExists(req.Email)
	if err != nil {
		api.HandleError(g, api.NewInternalServerError(err))
		return
	}
	if exists {
		api.HandleError(g, api.EmailExists())
		return
	}

	err = database.CreateUser(req.AccountType, &req.Avatar, req.DisplayName, req.Email, &req.ExternalID, req.FullName, req.Password, &req.Phone)
	if err != nil {
		api.HandleError(g, api.NewInternalServerError(err))
		return
	}

	result, err := database.GetUserByCreds(req.Email, req.Password)
	if err != nil {
		api.HandleError(g, api.NewInternalServerError(err))
		return
	}

	claims := utils.UserClaims{
		UserID:        result.IDString,
		Email:         result.Email,
		EmailVerified: result.EmailVerified,
	}

	secretKey := vault.Envars["TOKEN_SECRET"].(string)

	token, err := utils.GenerateToken(claims, []byte(secretKey))
	if err != nil {
		api.HandleError(g, api.OperationFailed())
		return
	}

	api.HandleSuccess(g, http.StatusOK, gin.H{
		"user":  result,
		"token": token,
	})
}
