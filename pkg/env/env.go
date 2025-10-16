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
	REDIS_URL    = ""
	DSN          = ""
	TOKEN_SECRET = ""
)

func InitalizeEnvs() {
	env := os.Getenv("ENV")

	if env == "" {
		logs.Warn("GetCurrentEnv", "current env is not set, defaulting to dev", nil)
		env = "dev"
	}

	switch env {
	case "dev":
		REDIS_URL = "redis:6379"
		GIN_MODE = "release"
	case "prod":
		REDIS_URL = "redis:6379"
		GIN_MODE = "release"
	default:
		REDIS_URL = "localhost:6379"
		GIN_MODE = "release"
	}
}
