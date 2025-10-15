package vault

import (
	"auth-service/pkg/logs"
	"log"
	"log/slog"
	"os"

	"auth-service/pkg/env"

	vaultclient "github.com/hashicorp/vault/api"
	"github.com/joho/godotenv"
)

type VaultOptions struct {
	DefaultToEnv bool
}

func InitializeVault() {
	vaultEnv := os.Getenv("VAULT_ENV")
	vaultAddress := os.Getenv("VAULT_ADDRESS")
	vaultToken := os.Getenv("VAULT_TOKEN")

	if vaultEnv == "local" {
		LoadEnvars()
		return
	}

	client, err := vaultclient.NewClient(&vaultclient.Config{
		Address: vaultAddress,
	})
	if err != nil {
		logs.Error("InitializeVault", "Error initializing Vault client: %s.", err)
		LoadEnvars()
		return
	}

	resp, err := client.Sys().Health()
	if err != nil {
		logs.Error("InitializeVault", "Error checking Vault health: %s. Falling back to .env variables.", err)
		LoadEnvars()
		return
	}

	if resp != nil && resp.Sealed {
		log.Println("Vault is sealed. Cannot proceed with fetching secrets.")
		LoadEnvars()
		return
	}

	if resp != nil && resp.Initialized && resp.Standby {
		log.Printf("Vault server is reachable but potentially has access issues.")
		LoadEnvars()
		return
	}

	client.SetToken(vaultToken)
	path := "pickside/data/credentials"

	secret, err := client.Logical().Read(path)
	if err != nil {
		log.Fatalf("Error reading secret: %s", err)
	}

	if secret == nil || secret.Data == nil {
		log.Fatal("Secret data is nil. Ensure the secret path is correct and the secret exists.")
	}

	data, ok := secret.Data["data"].(map[string]any)
	if !ok {
		log.Fatal("Secret structure is not as expected. Unable to find 'data' map.")
	}

	slog.Info("Connected to Vault", "success", true)

	env.DSN = data["DSN"].(string)
	env.TOKEN_SECRET = data["TOKEN_SECRET"].(string)
}

func LoadEnvars() {
	err := godotenv.Load()
	if err != nil {
		os.Exit(0)
	}

	env.DSN = os.Getenv("DSN")
	env.TOKEN_SECRET = os.Getenv("TOKEN_SECRET")

	logs.Info("LoadEnvars", "Loaded environment variables from .env file.", nil)
}
