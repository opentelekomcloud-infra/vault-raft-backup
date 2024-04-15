package app

import (
	"log"

	"github.com/opentelekomcloud-infra/vault-raft-backup/pkg/obs"
	"github.com/opentelekomcloud-infra/vault-raft-backup/pkg/vault"
)

// Backup creates a Raft snapshot of Vault data and uploads it to OBS.
// It handles all clients initializations, Vault login, snapshot creation, and OBS upload.
func Backup(vaultConfig *vault.Config, obsConfig *obs.Config) (err error) {
	obsClient, err := obs.InitClient()
	if err != nil {
		log.Fatalf("Failed to create OBS client: %s\n", err)
	}

	vaultClient, err := vault.InitClient(vaultConfig)
	if err != nil {
		log.Fatalf("Failed to create Vault client: %s\n", err)
	}

	if err := vaultClient.LoginAppRole(vaultConfig.RoleID, vaultConfig.SecretID); err != nil {
		log.Fatalf("Error logging into Vault: %s\n", err)
	}

	snapshotBuffer, err := vaultClient.CreateRaftSnapshot()
	if err != nil {
		log.Fatalf("Failed to create Raft snapshot: %s", err)
	}

	if err := obsClient.UploadToOBS(snapshotBuffer, obsConfig); err != nil {
		log.Fatalf("Failed to upload snapshot to OBS: %s", err)
	}
	return nil
}
