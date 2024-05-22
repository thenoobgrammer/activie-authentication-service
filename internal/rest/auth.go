package rest

import (
	"auth-service/internal/service"
	"auth-service/pkg/models"
	"auth-service/pkg/utils"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Me(g *gin.Context) {
	bearer := utils.GetBearer(g)

	result, err := service.QueryUserFromBearer(*bearer)
	if err != nil {
		g.JSON(http.StatusForbidden, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong",
			"success": false,
		})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"message": "User successfully matched",
		"result":  result,
		"success": true,
	})
}
func Authenticate(g *gin.Context) {
	var req models.AuthenticationRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": errors.New("bad request").Error()})
		return
	}

	service.HandleAuthentication(req)

	result, err := service.QueryUserFromCredentials(req)
	if err != nil {
		g.JSON(http.StatusNotFound, gin.H{
			"error":   err.Error(),
			"message": "Email or password do not match",
			"success": false,
		})
		return
	}

	bearer, err := service.IssueToken(models.UserClaims{
		UserID:        result.IDString,
		Email:         result.Email,
		EmailVerified: result.EmailVerified,
	})
	if err != nil {
		g.JSON(http.StatusNotFound, gin.H{
			"error":   err.Error(),
			"message": "Problem issuing token",
			"success": false,
		})
		return
	}

	g.Header("Authorization", fmt.Sprintf("Bearer %s", bearer))

	g.JSON(http.StatusOK, gin.H{
		"message": "Credentials match",
		"result":  result,
		"token":   bearer,
		"success": true,
	})
}
func Signup(g *gin.Context) {
	var req models.CreateUserRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	req.AccountType = "system"

	result, err := service.CreateUser(req)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong",
			"success": false,
		})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"message": "Account created",
		"result":  result,
		"success": true,
	})
}
