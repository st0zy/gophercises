/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/st0zy/gophercises/secret/secret"
)

var vault *secret.FileVault
var vaultDb string
var encryptionKey string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "secret",
	Short: "Store and Fetch secrets",
	Long: `Secret is, as you might have guessed, a CLI to store and fetch secrets.
	Secrets are stored to a file called secrets.db in the current path.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	//Flags aren't working currently. Need to be fixed
	// PS : flags are set only after rootCmd.Execute() is called
	rootCmd.PersistentFlags().StringVar(&vaultDb, "file", "secrets.db", "Vault DB path")
	rootCmd.PersistentFlags().StringVar(&encryptionKey, "key", "test123", "Key used to encrypt/decrypt")
	vault = secret.NewFileVault(secret.WithVaultPath(vaultDb), secret.WithEncryptionKey(encryptionKey))

}
