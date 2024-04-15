package cmd

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Global variables for command-line flags.
var (
	vaultRoleID   string
	vaultSecretID string
	vaultAddr     string
	vaultTimeout  time.Duration
	obsBucketName string
	obsObjectName string
	accessKey     string
	secretKey     string
	authURL       string
	domainName    string
	projectName   string
)

// rootCmd defines the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "vault-raft-backup",
	Short: "Backup Vault data using Raft snapshots",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			log.Fatalf("failed to show help: %s\n", err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

// initConfig configures the application by binding environment variables and setting up the configuration.
func initConfig() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	if err := viper.BindPFlags(backupCmd.Flags()); err != nil {
		log.Fatalf("cannot bind Viper to flags:%v", err)
	}
	checkAndSetEnv()
}

func init() {
	cobra.OnInitialize(initConfig) // Initialize the configuration when the Cobra application starts.

	// Define command line flags and mark required flags.
	backupCmd.Flags().StringVarP(&vaultRoleID, "vault-role-id", "", "", "Vault AppRole role ID (required)")
	if err := backupCmd.MarkFlagRequired("vault-role-id"); err != nil {
		log.Fatalf("error on marking flag as required: %v", err)
	}
	backupCmd.Flags().StringVarP(&vaultSecretID, "vault-secret-id", "", "", "Vault AppRole secret ID (required)")
	if err := backupCmd.MarkFlagRequired("vault-secret-id"); err != nil {
		log.Fatalf("error on marking flag as required: %v", err)
	}
	backupCmd.Flags().StringVarP(&vaultAddr, "vault-address", "", "https://127.0.0.1:8200", "vault address")
	backupCmd.PersistentFlags().DurationVar(&vaultTimeout, "vault-timeout", 60*time.Second, "vault client timeout")
	backupCmd.Flags().StringVarP(&obsBucketName, "obs-bucket-name", "", "", "OBS bucket name (required)")
	if err := backupCmd.MarkFlagRequired("obs-bucket-name"); err != nil {
		log.Fatalf("error on marking flag as required: %v", err)
	}
	backupCmd.Flags().StringVarP(&obsObjectName, "obs-object-name", "", "vault-raft-backup.snap", "OBS object name")
	backupCmd.Flags().StringVarP(&accessKey, "os-access-key", "", "", "Access Key for authentication (required)")
	if err := backupCmd.MarkFlagRequired("os-access-key"); err != nil {
		log.Fatalf("error on marking flag as required: %v", err)
	}
	backupCmd.Flags().StringVarP(&secretKey, "os-secret-key", "", "", "Secret Key for authentication (required)")
	if err := backupCmd.MarkFlagRequired("os-secret-key"); err != nil {
		log.Fatalf("error on marking flag as required: %v", err)
	}
	backupCmd.Flags().StringVarP(&authURL, "os-auth-url", "", "https://iam.eu-de.otc.t-systems.com/v3", "Authentication URL")
	backupCmd.Flags().StringVarP(&domainName, "os-domain-name", "", "", "Domain name (required)")
	if err := backupCmd.MarkFlagRequired("os-domain-name"); err != nil {
		log.Fatalf("error on marking flag as required: %v", err)
	}
	backupCmd.Flags().StringVarP(&projectName, "os-project-name", "", "eu-de", "Project name")
}

// checkAndSetEnv ensures all required environment variables are set, providing defaults if necessary.
func checkAndSetEnv() {
	envVariables := map[string]string{
		"OS_ACCESS_KEY":   accessKey,
		"OS_SECRET_KEY":   secretKey,
		"OS_AUTH_URL":     authURL,
		"OS_DOMAIN_NAME":  domainName,
		"OS_PROJECT_NAME": projectName,
	}

	for key, value := range envVariables {
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}
