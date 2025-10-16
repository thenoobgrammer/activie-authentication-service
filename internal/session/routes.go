package session

import (
	"auth-service/internal/token"
	"auth-service/internal/user"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func AttachHandlers(rg *gin.RouterGroup, db *sql.DB) {
	repo := NewRepository(db)
	userRepo := user.NewRepository(db)
	tokenRepo := token.NewRepository(db)
	service := NewService(repo, userRepo, tokenRepo)
	handler := NewHandler(service)

	rg.POST("/token", handler.IssueToken)
	rg.POST("/token/validate", handler.ValidateToken)
	rg.POST("/token/invalidate", handler.InvalidateToken)
	rg.POST("/token/refresh", handler.RefreshToken)
	rg.POST("/anonymous/session", handler.InitiateAnonymousSession)
}
