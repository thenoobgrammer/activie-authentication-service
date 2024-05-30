package main

import (
	"auth-service/internal/database"
	"auth-service/internal/redis"
	"auth-service/internal/rest"
	"auth-service/internal/vault"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	dsn := vault.Envars["DSN"].(string)
	database.Initialize(dsn)
	defer database.Close()

	redis.InitializeClient()

	gin.SetMode(os.Getenv("GIN_MODE"))
	g := gin.Default()
	g.Use(cors.New(buildCors()))

	//Health
	g.GET("/health", rest.GetHealth)

	//Session
	g.POST("/session/check", rest.CheckSession)
	g.POST("/session/new", rest.NewSession)
	g.POST("/session/end", rest.EndSession)

	//Token
	g.POST("/token/request", rest.RequestToken)
	g.POST("/token/refresh", rest.RefreshToken)

	PrintServiceInformation()

	if err := g.Run(":8081"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}

func PrintServiceInformation() {
	log.Printf("Mode %s", os.Getenv("GIN_MODE"))
	log.Printf("Service name: %s", os.Getenv("SERVICE_NAME"))
	log.Printf("Version: %s", os.Getenv("SERVICE_VERSION"))
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
