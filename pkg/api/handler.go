package api

import (
	"auth-service/pkg/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleError(g *gin.Context, err error) {
	switch e := err.(type) {
	case APIError:
		g.JSON(e.StatusCode, e)
	case InternalServerError:
		utils.LogError(g.FullPath(), "internal server error", err)
		apiErr := NewAPIError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
		g.JSON(apiErr.StatusCode, apiErr)
	default:
		utils.LogError(g.FullPath(), "internal server error", err)
		apiErr := NewAPIError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
		g.JSON(apiErr.StatusCode, apiErr)
	}
}

func HandleSuccess(g *gin.Context, statusCode int, data interface{}) {
	g.JSON(statusCode, data)
}
