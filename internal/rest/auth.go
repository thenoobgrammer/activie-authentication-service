package rest

import (
	"auth-service/internal/service"
	"auth-service/pkg/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthenticateBearer(g *gin.Context) {
	bearer := g.Query("bearerToken")

	if bearer == "" {
		g.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "Unsupported authentication method",
			"success": false,
		})
		return
	}

	ok, err := service.ValidateFromBearer(bearer)
	if err != nil {
		g.JSON(http.StatusNotFound, gin.H{
			"authenticated": false,
			"error":         err.Error(),
			"message":       "Email or password do not match",
			"success":       false,
		})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"authenticated": ok,
		"message":       "Bearer valid",
		"success":       true,
	})
}

func AuthenticateCredentials(g *gin.Context) {
	var req models.SystemUserCrendetialsRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": errors.New("bad request").Error()})
		return
	}

	ok, err := service.ValidateCredentials(req)
	if err != nil {
		g.JSON(http.StatusNotFound, gin.H{
			"authenticated": false,
			"error":         err.Error(),
			"message":       "Email or password do not match",
			"success":       false,
		})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"authenticated": ok,
		"message":       "Credentials match",
		"success":       true,
	})
}
