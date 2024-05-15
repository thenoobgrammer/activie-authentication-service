package rest

import (
	"auth-service/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateSession(g *gin.Context) {
	g.JSON(http.StatusNotImplemented, gin.H{"message": "Method not implemented"})
}
func TimeoutSession(g *gin.Context) {
	g.JSON(http.StatusNotImplemented, gin.H{"message": "Method not implemented"})
}
func Logout(g *gin.Context) {
	utils.ClearBearer(g)

	g.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
		"success": true,
	})
}
