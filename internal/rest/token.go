package rest

import (
	"auth-service/internal/service"
	"auth-service/pkg/models"
	"auth-service/pkg/utils"
	"log"
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
func VerifyToken(g *gin.Context) {
	bearerToken := utils.GetBearer(g)

	err := service.ValidateToken(*bearerToken)
	if err != nil {
		log.Println("Invalid token", err.Error())
		g.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   err.Error(),
			"message": "Unauthorized",
			"success": false,
		})
		return
	}

	if service.TokenExpired(*bearerToken) {
		refreshedBearer, err := service.RefreshToken(*bearerToken)
		if err != nil {
			g.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   err.Error(),
				"message": "Error while refreshing token, please relogin",
				"success": false,
			})
			return
		}
		utils.SetBearer(g, *refreshedBearer)
	}

	claims, err := service.GetClaims(*bearerToken)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Error while fetching claims",
			"success": false,
		})
	}

	g.JSON(http.StatusOK, claims)
}
