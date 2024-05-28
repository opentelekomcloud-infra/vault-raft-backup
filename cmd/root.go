package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Global variables to hold command line flags for configuration.
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
	cfgFile       string
)

// rootCmd defines the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "vault-raft-backup",
	Short: "Backup Vault data using Raft snapshots",
	Run: func(cmd *cobra.Command, args []string) {
		// Show help if no subcommands are called.
		if err := cmd.Help(); err != nil {
			log.Fatalf("Failed to show help: %s\n", err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error during command execution: %s\n", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig) // Initialize the configuration when the Cobra application starts.

	// Define flags and configuration settings.
	rootCmd.PersistentFlags().StringVarP(&obsBucketName, "obs-bucket-name", "", "", "OBS bucket name (required)")
	rootCmd.PersistentFlags().StringVarP(&obsObjectName, "obs-object-name", "", "vault-raft-backup.snap", "OBS object name")
	rootCmd.PersistentFlags().StringVarP(&accessKey, "os-access-key", "", "", "OTC Access Key for authentication (required)")
	rootCmd.PersistentFlags().StringVarP(&secretKey, "os-secret-key", "", "", "OTC Secret Key for authentication (required)")
	rootCmd.PersistentFlags().StringVarP(&authURL, "os-auth-url", "", "https://iam.eu-de.otc.t-systems.com/v3", "OTC Authentication URL")
	rootCmd.PersistentFlags().StringVarP(&domainName, "os-domain-name", "", "", "OTC Domain name (required)")
	rootCmd.PersistentFlags().StringVarP(&projectName, "os-project-name", "", "eu-de", "OTC Project name")
	rootCmd.PersistentFlags().StringVarP(&vaultAddr, "vault-address", "", "https://127.0.0.1:8200", "Vault address")
	rootCmd.PersistentFlags().StringVarP(&vaultRoleID, "vault-role-id", "", "", "Vault AppRole role ID (required)")
	rootCmd.PersistentFlags().StringVarP(&vaultSecretID, "vault-secret-id", "", "", "Vault AppRole secret ID (required)")
	rootCmd.PersistentFlags().DurationVar(&vaultTimeout, "vault-timeout", 60*time.Second, "Vault client timeout")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vault-raft-backup.yaml)")

	// Mark required flags.
	markRequiredFlags()
}

// markRequiredFlags ensures that essential flags are set before the application runs.
func markRequiredFlags() {
	requiredFlags := []string{"vault-role-id", "vault-secret-id", "obs-bucket-name", "os-access-key", "os-secret-key", "os-domain-name"}
	for _, flag := range requiredFlags {
		if err := rootCmd.MarkPersistentFlagRequired(flag); err != nil {
			log.Fatalf("Error marking flag %s as required: %v", flag, err)
		}
	}
}

// initConfig configures the application by binding environment variables and setting up the configuration.
func initConfig() {
	v := viper.New()

	// Set the configuration file if provided
	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		// Find the home directory
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Error finding home directory: %v", err)
		}

		// Search for the config file in the home directory
		v.AddConfigPath(home)
		v.SetConfigType("yaml")
		v.SetConfigName(".vault-raft-backup")
	}

	// Read in the config file if it exists
	if err := v.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", v.ConfigFileUsed())
	}

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Bind flags to viper configuration
	bindFlags(rootCmd, v)
	// Ensure environment variables are set
	checkAndSetEnv()
}

// bindFlags binds each cobra flag to its associated viper configuration (environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			if err := cmd.PersistentFlags().Set(f.Name, fmt.Sprintf("%v", val)); err != nil {
				log.Fatalf("Error setting flag %s: %v", f.Name, err)
			}
		}
	})
}

// checkAndSetEnv ensures all required environment variables for auth package are set, providing defaults if necessary.
func checkAndSetEnv() {
	envVariables := map[string]string{
		"OS_ACCESS_KEY":   accessKey,
		"OS_SECRET_KEY":   secretKey,
		"OS_AUTH_URL":     authURL,
		"OS_DOMAIN_NAME":  domainName,
		"OS_PROJECT_NAME": projectName,
	}

	for key, value := range envVariables {
		if os.Getenv(key) == "" && value != "" {
			os.Setenv(key, value)
		}
	}
}
