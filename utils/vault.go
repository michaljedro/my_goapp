package utils

import (
    "log"
    vault "github.com/hashicorp/vault/api"
)

func NewVaultClient() (*vault.Client, error) {
    config := vault.DefaultConfig()
    config.Address = "http://vault:8200"

    client, err := vault.NewClient(config)
    if err != nil {
        return nil, err
    }

    log.Println("Connected to Vault")
    return client, nil
}
