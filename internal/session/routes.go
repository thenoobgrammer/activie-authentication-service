package session

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func AttachHandlers(rg *gin.RouterGroup, db *sql.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	rg.GET("/sessions", handler.GetActiveSession)
	rg.POST("/sessions/start", handler.StartSession)
	rg.DELETE("/sessions/end", handler.EndActiveSession)
}
