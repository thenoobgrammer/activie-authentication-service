package middlewares

import (
	"auth-service/pkg/api"

	"github.com/gin-gonic/gin"
)

const (
	Client    = "ZuI8XOfWzP" // For any users
	DevClient = "4CEAB"      // For dev testing purposes
)

func IsValidClient() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientId := c.GetHeader("Client-Id")
		if clientId == Client {
			c.Set("ClientName", "authenticated-client")
		} else if clientId == DevClient {
			c.Set("ClientName", "dev-client")
		} else {
			c.Set("ClientApproved", false)
			api.HandleError(c, api.InvalidHeaderRequest())
			return
		}
		c.Set("ClientApproved", true)
		c.Next()
	}
}
