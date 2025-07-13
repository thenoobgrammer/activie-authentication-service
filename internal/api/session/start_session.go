package api_session

import (
	database_session "auth-service/internal/database/session"
	"auth-service/internal/models"
	"auth-service/pkg/api"
	"auth-service/pkg/env"
	"auth-service/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"
)

func StartSession(g *gin.Context) {
	var req models.SessionRequirements

	if err := g.ShouldBindJSON(&req); err != nil {
		pretty.Println("HERE", req)
		api.HandleError(g, api.InvalidJSON())
		return
	}

	key := env.TOKEN_SECRET

	token := utils.GenerateToken(utils.UserClaims{
		AccountType: req.AccountType,
		Email:       req.UserEmail,
		UserID:      req.UserId,
		UserRoles:   req.UserRoles,
	}, []byte(key))
	if token == nil {
		api.HandleError(g, api.FailedToGenerateBearerToken())
		return
	}

	database_session.DeleteFromUserId(req.UserId) // delete session if one already exists

	sessionId := database_session.Create(req, *token)
	if sessionId == nil {
		api.HandleError(g, api.FailedToCreateSession())
		return
	}

	ctxSession := database_session.Retrieve(*sessionId)

	api.HandleSuccess(g, http.StatusOK, gin.H{
		"ctxSession": ctxSession,
		"token":      token,
	})
}
