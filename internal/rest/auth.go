package rest

import (
	"auth-service/internal/service"
	"auth-service/pkg/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthenticateBearer(g *gin.Context) {
	bearer := g.Query("token")

	if bearer == "" {
		g.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "Unsupported authentication method",
			"success": false,
		})
		return
	}

	authenticated := service.ValidateFromBearer(bearer)
	if !authenticated {
		g.JSON(http.StatusUnauthorized, gin.H{
			"authenticated": false,
			"message":       "Email or password do not match",
			"success":       false,
		})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"authenticated": authenticated,
		"success":       true,
	})
}
func AuthenticateCredentials(g *gin.Context) {
	var req models.SystemUserCrendetialsRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": errors.New("bad request").Error()})
		return
	}

	authenticated := service.ValidateCredentials(req)
	if !authenticated {
		g.JSON(http.StatusUnauthorized, gin.H{
			"authenticated": false,
			"message":       "Email or password do not match",
			"success":       false,
		})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"authenticated": authenticated,
		"message":       "Credentials match",
		"success":       true,
	})
}
