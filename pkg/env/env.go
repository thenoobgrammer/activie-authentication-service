package env

import (
	"auth-service/pkg/logs"
	"os"
)

const (
	SERVICE_NAME = "Activie authentication service"
	VERSION      = "0.1.1"
)

var (
	GIN_MODE     = ""
	DSN          = os.Getenv("DSN")
	TOKEN_SECRET = os.Getenv("TOKEN_SECRET")
)

func InitalizeEnvs() {
	env := os.Getenv("ENV")

	if env == "" {
		logs.Warn("GetCurrentEnv", "current env is not set, defaulting to dev", nil)
		env = "local"
	}

	switch env {
	case "local":
		GIN_MODE = "debug"
	case "prod":
		GIN_MODE = "release"
	default:
		GIN_MODE = "release"
	}
}
