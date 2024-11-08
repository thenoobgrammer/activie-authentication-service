package utils

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func ExtractToken(authHeader string) *string {
	if authHeader == "" {
		return nil
	}
	slice := strings.Split(authHeader, " ")
	if len(slice) < 2 {
		return nil
	}
	if slice[0] == "" || slice[0] != "Bearer" {
		return nil
	}
	if slice[1] == "" {
		return nil
	}
	return &slice[1]
}

func SetBearer(g *gin.Context, token string) {
	g.Header("Authorization", fmt.Sprintf("Bearer %s", token))
}

func ClearBearer(g *gin.Context) {
	g.Request.Header.Del("Authorization")
}
