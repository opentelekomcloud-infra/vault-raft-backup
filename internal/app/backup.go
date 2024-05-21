package app

import (
	"log"

	"github.com/opentelekomcloud-infra/vault-raft-backup/pkg/obs"
	"github.com/opentelekomcloud-infra/vault-raft-backup/pkg/vault"
)

// Backup creates a Raft snapshot of Vault data and uploads it to OBS.
// It handles all clients initializations, Vault login, snapshot creation, and OBS upload.
func Backup(vaultConfig *vault.Config, obsConfig *obs.Config) error {
	log.Println("Starting OBS client initialization...")
	obsClient, err := obs.InitClient()
	if err != nil {
		log.Fatalf("Failed to create OBS client: %v", err)
	}
	log.Println("OBS client initialized successfully...")

	log.Println("Starting Vault client initialization...")
	vaultClient, err := vault.InitClient(vaultConfig)
	if err != nil {
		log.Fatalf("Failed to create Vault client: %v", err)
	}
	log.Println("Vault client initialized successfully...")

	log.Println("Logging into Vault...")
	if err := vaultClient.LoginAppRole(vaultConfig.RoleID, vaultConfig.SecretID); err != nil {
		log.Fatalf("Error logging into Vault: %v", err)
	}
	log.Println("Logged into Vault successfully...")

	log.Println("Creating Raft snapshot...")
	snapshotBuffer, err := vaultClient.CreateRaftSnapshot()
	if err != nil {
		log.Fatalf("Failed to create Raft snapshot: %v", err)
	}
	log.Println("Raft snapshot created successfully...")

	log.Println("Uploading snapshot to OBS...")
	if err := obsClient.UploadToOBS(snapshotBuffer, obsConfig); err != nil {
		log.Fatalf("Failed to upload snapshot to OBS: %v", err)
	}
	log.Println("Snapshot uploaded to OBS successfully...")

	return nil
}
