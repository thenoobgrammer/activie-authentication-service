package api

import (
	"auth-service/pkg/api"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Logout(g *gin.Context) {
	authHeader := g.GetHeader("Authorization")
	if authHeader == "" {
		api.HandleError(g, api.InvalidBearer())
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		api.HandleError(g, api.InvalidBearer())
		return
	}

	api.HandleSuccess(g, http.StatusOK, "successfully logged out")
}
