/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/st0zy/gophercises/task/storage"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		taskId, err := cmd.Flags().GetUint64("taskId")

		if err != nil {
			fmt.Println("failed to parse the taskId with err ", err)
		}

		err = storage.DoTask(taskId)
		if err != nil {
			fmt.Println("failed to mark task as completed.", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(doCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// doCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	doCmd.Flags().Uint64P("taskId", "t", 0, "taskId to delete")
	doCmd.MarkFlagRequired("taskId")
}
