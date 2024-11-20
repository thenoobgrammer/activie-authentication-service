package env

import (
	"auth-service/pkg/utils"
	"fmt"
	"os"
)

var (
	GIN_MODE      = ""
	VAULT_ADDRESS = ""
	VAULT_TOKEN   = ""
	SERVICE_NAME  = "Authentication service"
	VERSION       = ""

	DSN          = ""
	TOKEN_SECRET = ""
)

func InitalizeEnvs() {
	env := os.Getenv("ENV")

	if env == "" {
		utils.LogWarn("GetCurrentEnv", "current env is not set, defaulting to dev", nil)
		env = "dev"
	}

	version, err := os.ReadFile("./.version")
	if err != nil {
		utils.LogError("GetCurrentEnv", "error reading .version file. make sure it's read, otherwise version will be set to 'unknown", err)
		VERSION = "unknown"
	} else {
		VERSION = fmt.Sprint(string(version))
	}

	switch env {
	case "prod":
		GIN_MODE = "release"
	default:
		GIN_MODE = "debug"
	}
}
