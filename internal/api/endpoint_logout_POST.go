package api

import (
	api "auth-service/internal/api/msgs"
	"auth-service/internal/redis"
	"auth-service/pkg/utils"
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

	if err := invalidateTokenInCache(tokenString); err != nil {
		api.HandleError(g, api.InvalidBearer())
		return
	}

	api.HandleSuccess(g, http.StatusOK, "successfully logged out")
}

func invalidateTokenInCache(token string) error {
	err := redis.InvalidateActiveSession(token)
	if err != nil {
		utils.LogError("invalidateTokenInCache", "Failed to invalidate token", err)
		return err
	}
	return nil
}
