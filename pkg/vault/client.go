package vault

import (
	"bytes"
	"fmt"
	"time"

	vault "github.com/hashicorp/vault/api"
)

// Client wraps the HashiCorp Vault API client with additional functionality.
type Client struct {
	vaultClient *vault.Client
}

// Config defines the necessary configuration to initialize a Vault client.
type Config struct {
	Address  string
	Timeout  time.Duration
	RoleID   string
	SecretID string
}

// InitClient initializes and returns a new Vault client using the provided configuration.
func InitClient(config *Config) (*Client, error) {
	vaultConfig := vault.DefaultConfig()
	vaultConfig.Address = config.Address

	client, err := vault.NewClient(vaultConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating Vault client: %w", err)
	}

	return &Client{
		vaultClient: client,
	}, nil
}

// LoginAppRole authenticates to Vault using the AppRole method with the provided roleID and secretID.
func (c *Client) LoginAppRole(roleID, secretID string) error {
	data := map[string]interface{}{
		"role_id":   roleID,   // Role identifier for AppRole.
		"secret_id": secretID, // Secret corresponding to the Role ID.
	}
	secret, err := c.vaultClient.Logical().Write("auth/approle/login", data)
	if err != nil {
		return fmt.Errorf("error logging in with AppRole: %w", err)
	}
	c.vaultClient.SetToken(secret.Auth.ClientToken)
	return nil
}

// CreateRaftSnapshot creates a snapshot of the Vault's Raft storage backend.
// Returns a buffer containing the snapshot data.
func (c *Client) CreateRaftSnapshot() (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	if err := c.vaultClient.Sys().RaftSnapshot(buf); err != nil {
		return nil, fmt.Errorf("error creating Raft snapshot: %w", err)
	}
	return buf, nil
}
