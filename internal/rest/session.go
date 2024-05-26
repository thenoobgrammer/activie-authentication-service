package rest

import (
	"auth-service/internal/redis"
	"auth-service/internal/service"
	"auth-service/pkg/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewSession(g *gin.Context) {
	var req models.SystemUserCrendetialsRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.New("bad request").Error()})
		return
	}

	authenticated, emailVerified, userID := service.ValidateCredentials(req)
	if !authenticated {
		g.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"authenticated": false,
			"message":       "Email or password do not match",
			"success":       false,
		})
		return
	}

	token, err := service.RequestToken(models.UserClaims{
		Email:         req.Email,
		EmailVerified: emailVerified,
		UserID:        *userID,
	})
	if err != nil {
		g.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"authenticated": false,
			"message":       "Something went wrong, please try again",
			"success":       false,
		})
		return
	}

	g.JSON(http.StatusOK, token)
}
func CheckSession(g *gin.Context) {
	bearer := g.Query("token")

	if bearer == "" {
		g.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "Unsupported authentication method",
			"session": false,
		})
		return
	}

	authenticated, claims := service.ValidateFromBearer(bearer)
	if !authenticated {
		g.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"authenticated": false,
			"message":       "Email or password do not match",
			"session":       false,
		})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"email":         claims.Email,
		"emailVerified": claims.EmailVerified,
		"userId":        claims.UserID,
		"session":       true,
	})
}
func EndSession(g *gin.Context) {
	bearer := g.Query("token")

	redis.DeleteSession(bearer)
	redis.BlacklistSession(bearer)

	g.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
		"success": true,
	})
}
