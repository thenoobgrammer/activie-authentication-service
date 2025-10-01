package session

import (
	"auth-service/pkg/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) StartSession(g *gin.Context) {
	var params NewSessionParams

	if err := g.ShouldBindJSON(&params); err != nil {
		api.HandleError(g, api.InvalidJSON())
		return
	}

	activeSession, token, err := h.service.StartAuthenticatedSession(params)
	if err != nil || activeSession == nil || token == nil {
		api.HandleError(g, api.FailedToCreateSession())
		return
	}

	api.HandleSuccess(g, http.StatusOK, *token)
}

func (h *Handler) StartAnonymousSession(c *gin.Context) {
	var params AnonymousNewSessionParams

	if err := c.ShouldBindJSON(&params); err != nil {
		api.HandleError(c, api.InvalidJSON())
		return
	}

	params.ClientID = c.GetHeader("X-Client-ID")
	params.IP = c.Request.RemoteAddr

	anonymousSession, token, err := h.service.StartAnonymousSession(params)
	if err != nil || anonymousSession == nil || token == nil {
		api.HandleError(c, api.FailedToCreateSession())
		return
	}

	api.HandleSuccess(c, http.StatusOK, *token)
}

// 9#z95W7ho5uoGjcN
func (h *Handler) GetActiveSession(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		api.HandleError(c, api.InvalidBearer())
		return
	}

	session := h.service.GetActiveSession(token)
	if session == nil {
		api.HandleError(c, api.SessionNotFound())
		return
	}

	api.HandleSuccess(c, http.StatusOK, session)
}

func (h *Handler) EndActiveSession(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		api.HandleError(c, api.InvalidBearer())
		return
	}

	ok := h.service.EndActiveSession(token)
	if !ok {
		api.HandleError(c, api.FailedToDeleteSession())
		return
	}

	api.HandleSuccess(c, http.StatusOK, true)
}
