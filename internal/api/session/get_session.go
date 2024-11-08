package api_session

import (
	database_session "auth-service/internal/database/session"
	"auth-service/pkg/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSession(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		api.HandleError(c, api.InvalidBearer())
		return
	}

	session := database_session.RetrieveFromToken(token)
	if session == nil {
		api.HandleError(c, api.SessionNotFound())
		return
	}

	api.HandleSuccess(c, http.StatusOK, session)
}
