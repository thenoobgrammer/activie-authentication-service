package health

import (
	"github.com/gin-gonic/gin"
)

func AttachHandlers(rg *gin.RouterGroup) {
	handler := NewHandler()
	
	rg.GET("/health", handler.GetHealth)
}
