package api_session

import (
	database_session "auth-service/internal/database/session"
	"auth-service/pkg/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func EndSession(g *gin.Context) {
	token := g.Query("token")
	if token == "" {
		api.HandleError(g, api.SessionNotFound())
		return
	}

	isDeleted := database_session.DeleteFromToken(token)

	api.HandleSuccess(g, http.StatusOK, isDeleted)
}
