package session

import (
	"auth-service/pkg/api"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) IssueToken(c *gin.Context) {
	var params ClaimRequirementsParams

	if err := c.ShouldBindJSON(&params); err != nil {
		api.AuthServiceReponse(c, http.StatusBadRequest, "invalid.token.params")
		return
	}

	token, err := h.service.GetTokenAndStartSession(params)
	if err != nil || token == nil {
		api.AuthServiceReponse(c, http.StatusInternalServerError, "failed.to.create.session")
		return
	}

	api.AuthServiceReponse(c, http.StatusOK, *token)
}

func (h *Handler) ValidateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		api.AuthServiceReponse(c, http.StatusUnauthorized, false)
		return
	}
	match := regexp.MustCompile(`^Bearer\s(.+)$`).FindStringSubmatch(authHeader)
	if len(match) != 2 {
		api.AuthServiceReponse(c, http.StatusUnauthorized, false)
		return
	}

	tokenStr := match[1]

	if err := h.service.ValidateToken(tokenStr); err != nil {
		api.AuthServiceReponse(c, http.StatusUnauthorized, false)
		return
	}

	api.AuthServiceReponse(c, http.StatusOK, true)
}

func (h *Handler) InvalidateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		api.AuthServiceReponse(c, http.StatusUnauthorized, false)
		return
	}
	match := regexp.MustCompile(`^Bearer\s(.+)$`).FindStringSubmatch(authHeader)
	if len(match) != 2 {
		api.AuthServiceReponse(c, http.StatusUnauthorized, false)
		return
	}

	tokenStr := match[1]

	if err := h.service.InvalidateToken(tokenStr); err != nil {
		api.AuthServiceReponse(c, http.StatusUnauthorized, false)
		return
	}

	api.AuthServiceReponse(c, http.StatusOK, true)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var params RefreshTokenParams
	if err := c.ShouldBindJSON(&params); err != nil {
		api.AuthServiceReponse(c, http.StatusBadRequest, "invalid.refresh.token")
		return
	}

	token, err := h.service.RefreshToken(params.RefreshToken)
	if err != nil {
		api.AuthServiceReponse(c, http.StatusInternalServerError, "failed.to.refresh.token")
		return
	}

	api.AuthServiceReponse(c, http.StatusOK, token)
}

func (h *Handler) InitiateAnonymousSession(c *gin.Context) {
	token, err := h.service.CreateAnonymousSession()
	if err != nil {
		api.AuthServiceReponse(c, http.StatusInternalServerError, "failed.to.initiate.anonymous.token")
		return
	}

	api.AuthServiceReponse(c, http.StatusOK, token)
}
