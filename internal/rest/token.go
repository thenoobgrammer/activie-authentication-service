package rest

import (
	"auth-service/internal/service"
	"auth-service/pkg/models"
	"auth-service/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IssueToken(g *gin.Context) {
	var queryParams models.UserClaims

	if err := g.ShouldBindQuery(&queryParams); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := service.IssueToken(queryParams)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusInternalServerError, "Error while issuing new token. Please reload and relogin.")
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
