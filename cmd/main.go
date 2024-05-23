package main

import (
	"auth-service/internal/database"
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

	gin.SetMode(os.Getenv("GIN_MODE"))
	g := gin.Default()
	g.Use(cors.New(buildCors()))

	as := g.Group("/auth-service")

	//Health
	as.GET("/health", rest.GetHealth)

	//User
	as.POST("/authenticate/bearer", rest.AuthenticateBearer)
	as.POST("/authenticate/credentials", rest.AuthenticateCredentials)

	//Token
	as.POST("/issue-token", rest.IssueToken)
	as.POST("/refresh-token", rest.RefreshToken)

	//Session
	as.GET("/logout", rest.Logout)

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
