package utils

import (
	"log"
	"os"
)

func IsSecure() bool {
	env := os.Getenv("GO_ENV")
	log.Println("env", env)
	switch env {
	case "prod":
		return true
	default:
		return false
	}
}
