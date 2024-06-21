package main

import (
	"auth-service/internal/api"
	"auth-service/internal/database"
	"auth-service/internal/vault"
	"auth-service/pkg/env"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Env intialization
	env.InitalizeEnvs()

	// Vault initialization
	vault.InitializeVault()

	// DB initialization
	database.InitializeDB(vault.Envars["DSN"].(string))
	defer database.Close()

	// Loggin initialization
	logHandler := slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelDebug}).WithAttrs([]slog.Attr{slog.String("service", "authentication")})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)
	slog.Info("service started")

	gin.SetMode(os.Getenv("GIN_MODE"))
	g := gin.Default()
	g.Use(cors.New(buildCors()))

	g.GET("/health", api.GetHealth)

	g.POST("/bearer-auth", api.BearerAuthentication)
	g.POST("/change-password", api.ChangePassword)
	g.POST("/login", api.Login)
	g.POST("/logout", api.Logout)
	g.POST("/signup", api.Signup)

	PrintServiceInformation()

	if err := g.Run(":8081"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func PrintServiceInformation() {
	log.Printf("Mode %s", env.GIN_MODE)
	log.Printf("Service name: %s", env.SERVICE_NAME)
	log.Printf("Version: %s", env.VERSION)
}

func buildCors() cors.Config {
	c := cors.DefaultConfig()
	c.AllowAllOrigins = false
	c.AllowCredentials = true
	c.AllowHeaders = []string{"Accept-Version", "Authorization", "Content-Type", "Origin", "X-Client-Version", "X-CSRF-Token", "X-Request-Id"}
	c.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	c.AllowWebSockets = false
	c.MaxAge = 24 * time.Hour

	c.AllowOriginFunc = func(origin string) bool {
		return true
	}
	return c
}
