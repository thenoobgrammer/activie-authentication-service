package vault

import (
	"auth-service/pkg/env"
	"log"
	"log/slog"

	vaultclient "github.com/hashicorp/vault/api"
)

var (
	Envars map[string]interface{}
)

func InitializeVault() {
	client, err := vaultclient.NewClient(&vaultclient.Config{
		Address: env.VAULT_ADDRESS,
	})
	if err != nil {
		log.Fatalf("Error initializing Vault client: %s", err)
	}

	vaultToken := env.VAULT_TOKEN
	client.SetToken(vaultToken)

	path := "pickside/data/credentials"

	secret, err := client.Logical().Read(path)
	if err != nil {
		log.Fatalf("Error reading secret: %s", err)
	}

	if secret == nil || secret.Data == nil {
		log.Fatal("Secret data is nil. Ensure the secret path is correct and the secret exists.")
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		log.Fatal("Secret structure is not as expected. Unable to find 'data' map.")
	}

	slog.Info("Connected to Vault", "success", true)

	Envars = data
}
