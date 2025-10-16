package api

import (
	"auth-service/pkg/logs"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleError(g *gin.Context, err error) {
	switch e := err.(type) {
	case APIError:
		g.JSON(e.StatusCode, e)
	case InternalServerError:
		logs.Error(g.FullPath(), "internal server error", err)
		apiErr := NewAPIError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
		g.JSON(apiErr.StatusCode, apiErr)
	default:
		logs.Error(g.FullPath(), "internal server error", err)
		apiErr := NewAPIError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
		g.JSON(apiErr.StatusCode, apiErr)
	}
}

func HandleSuccess(g *gin.Context, statusCode int, data interface{}) {
	g.JSON(statusCode, data)
}

func AuthServiceReponse(c *gin.Context, code int, data any) {
	c.JSON(code, data)
}
