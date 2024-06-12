package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WriteAPIError(g *gin.Context, apiErr APIError) {
	g.AbortWithStatusJSON(apiErr.StatusCode, apiErr)
}

func WriteAPISuccess(g *gin.Context, apiResponse APIResponse) {
	g.JSON(apiResponse.StatusCode, apiResponse)
}

func HandleError(g *gin.Context, err error) {
	switch e := err.(type) {
	case APIError:
		WriteAPIError(g, e)
	case InternalServerError:
		slog.Error("Internal server error", e.Err)
		//log.Printf("Internal server error: %v", e.Err) // Log the internal server error
		apiErr := NewAPIError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
		WriteAPIError(g, apiErr)
	default:
		slog.Error("Internal server error", err)
		//log.Printf("Unknown error: %v", err) // Log unknown errors
		apiErr := NewAPIError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
		WriteAPIError(g, apiErr)
	}
}

func HandleSuccess(g *gin.Context, statusCode int, data interface{}) {
	apiResponse := NewAPIResponse(statusCode, data)
	WriteAPISuccess(g, apiResponse)
}
