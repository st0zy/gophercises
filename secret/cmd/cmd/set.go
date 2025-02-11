/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Store a secret",
	Long:  `Store the secret key-value pair into the vault.`,
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]
		err := vault.Set(key, value)
		if err != nil {
			fmt.Println("failed to fetch the secret", err.Error())
		}
		fmt.Println("Stored successfully")
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
