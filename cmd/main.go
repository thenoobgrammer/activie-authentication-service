package main

import (
	"log"
	"log/slog"
	"os"
	"time"

	"auth-service/internal/api"
	api_session "auth-service/internal/api/session"
	"auth-service/internal/database"
	"auth-service/internal/vault"
	"auth-service/middlewares"
	"auth-service/pkg/env"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Env intialization
	env.InitalizeEnvs()

	// Vault initialization
	vault.InitializeVault()

	// DB initialization
	database.InitializeDB(env.DSN)
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
	g.Use(middlewares.RequestMetricsMiddleware())

	g.GET("/health", api.GetHealth)

	g.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Sessions
	g.GET("/sessions", api_session.GetSession)
	g.POST("/sessions/start", api_session.StartSession)
	g.DELETE("/sessions/end", api_session.EndSession)

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
