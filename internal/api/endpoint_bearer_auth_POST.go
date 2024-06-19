package api

import (
	"auth-service/internal/database"
	"auth-service/internal/vault"
	"auth-service/pkg/api"
	"auth-service/pkg/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func BearerAuthentication(g *gin.Context) {
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

	secretKey := vault.Envars["TOKEN_SECRET"].(string)

	claims, err := utils.GetClaims(tokenString, secretKey)
	if err != nil {
		api.HandleError(g, api.InvalidBearer())
		return
	}

	result, err := database.GetUserById(claims.UserID)
	if err != nil {
		api.HandleError(g, api.NewInternalServerError(err))
	} else if result == nil {
		api.HandleError(g, api.NotFound())
	}

	log.Println("result", result)
	api.HandleSuccess(g, http.StatusOK, gin.H{
		"user": result,
	})
}
