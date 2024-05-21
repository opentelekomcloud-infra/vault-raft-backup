package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/opentelekomcloud-infra/vault-raft-backup/internal/app"
	"github.com/opentelekomcloud-infra/vault-raft-backup/pkg/obs"
	"github.com/opentelekomcloud-infra/vault-raft-backup/pkg/vault"
)

// backupCmd defines a command for performing a backup of Vault data.
// The command uses Raft snapshots to capture the state of Vault's data and store it in OBS.
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Perform a backup of Vault data using Raft snapshot",
	Run: func(cmd *cobra.Command, args []string) {
		// Configuration for the OBS client.
		obsConfig := &obs.Config{
			BucketName: obsBucketName,
			ObjectName: obsObjectName,
		}

		// Configuration for the Vault client.
		vaultConfig := &vault.Config{
			RoleID:   vaultRoleID,
			SecretID: vaultSecretID,
			Address:  vaultAddr,
			Timeout:  vaultTimeout,
		}

		// Execute the backup process using the specified configurations.
		err := app.Backup(vaultConfig, obsConfig)
		if err != nil {
			log.Fatalf("Backup failed: %s\n", err)
		}
		log.Println("Backup completed successfully")
	},
}

func init() {
	rootCmd.AddCommand(backupCmd) // Register the backup command with the root command.
}
