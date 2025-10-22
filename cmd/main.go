package main

import (
	"auth-service/internal/health"
	"auth-service/internal/infra/database"
	"auth-service/internal/session"
	"auth-service/pkg/logs"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"auth-service/pkg/env"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	env.InitalizeEnvs()
	database.InitializeDB(env.DSN)

	gin.SetMode(env.GIN_MODE)
	engine := gin.Default()

	setupServer(engine)

	srv := http.Server{
		Addr:    ":8081",
		Handler: engine,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() { srv.ListenAndServe() }()

	<-quit

	cleanUpServices()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}
}

func printServiceInformation() {
	logs.Info("Mode: ", env.GIN_MODE)
	logs.Info("Service name: ", env.SERVICE_NAME)
	logs.Info("Version: ", env.VERSION)
}

func setupServer(engine *gin.Engine) {
	db := database.GetClient()
	if db == nil {
		panic("database connection failed")
	}

	engine.Use(cors.New(buildCors()))
	engine.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(204)
	})

	api := engine.Group("/")

	session.AttachHandlers(api, db)
	health.AttachHandlers(api)
}

func cleanUpServices() {
	database.Close()
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
