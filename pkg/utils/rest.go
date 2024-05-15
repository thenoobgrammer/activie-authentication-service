package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GetBearer(g *gin.Context) *string {
	authorizationHeader := strings.Split(g.GetHeader("Authorization"), " ")

	if len(authorizationHeader) < 2 || strings.ToLower(authorizationHeader[0]) != "bearer" {
		g.JSON(http.StatusForbidden, gin.H{
			"status":  "Unauthorized",
			"message": jwt.ErrTokenMalformed,
			"success": false,
		})
		return nil
	}

	return &authorizationHeader[1]
}

func SetBearer(g *gin.Context, token string) {
	g.Header("Authorization", fmt.Sprintf("Bearer %s", token))
}

func ClearBearer(g *gin.Context) {
	g.Request.Header.Del("Authorization")
}
