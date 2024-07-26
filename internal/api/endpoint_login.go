package api

import (
	"auth-service/internal/database"
	"auth-service/internal/vault"
	"auth-service/pkg/api"
	"auth-service/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (req *LoginRequest) Validate() map[string]string {
	var errors = make(map[string]string)

	if req.Email == "" {
		errors["email"] = "email is required"
	}
	if req.Password == "" {
		errors["password"] = "password is required"
	}

	return errors
}

func Login(g *gin.Context) {
	var req LoginRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		api.HandleError(g, api.InvalidJSON())
		return
	}

	if errors := req.Validate(); len(errors) > 0 {
		api.HandleError(g, api.InvalidRequestData(errors))
		return
	}

	result, err := database.GetUserByCreds(req.Email, req.Password)
	if err != nil {
		api.HandleError(g, api.InvalidCredentials())
		return
	}

	claims := utils.UserClaims{
		AccountType: result.AccountType,
		Email:       result.Email,
		Permissions: result.PermissionsString,
		Role:        result.Role,
		UserID:      result.IDString,
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
