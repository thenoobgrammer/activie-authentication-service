package rest

import (
	"auth-service/internal/models"
	"auth-service/internal/service"
	"auth-service/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequestToken(g *gin.Context) {
	var queryParams models.UserClaims

	if err := g.ShouldBindJSON(&queryParams); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := service.RequestToken(queryParams)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": true,
			"message": err.Error(),
		})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"status":  "Authorized",
		"message": "Token issuance",
		"result":  result,
		"success": true,
	})
}
func RefreshToken(g *gin.Context) {
	bearerToken := utils.GetBearer(g)

	result, err := service.RefreshToken(*bearerToken)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusInternalServerError, "Error while issuing new token. Please reload and relogin.")
		return
	}

	utils.SetBearer(g, *result)

	g.JSON(http.StatusOK, gin.H{
		"message": "Refreshed bearer",
		"success": true,
	})
}
