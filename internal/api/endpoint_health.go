package api

import (
	"auth-service/pkg/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHealth(g *gin.Context) {
	api.HandleSuccess(g, http.StatusOK, "service up")
}
