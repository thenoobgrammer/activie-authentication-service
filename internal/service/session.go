package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout(g *gin.Context) {
	g.Request.Header.Del("Authorization")

	g.JSON(http.StatusOK, gin.H{
		"message":     "Logged out successfully",
		"redirectUri": "/login",
		"success":     true,
	})
}
