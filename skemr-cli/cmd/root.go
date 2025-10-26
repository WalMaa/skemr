package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "skemr-cli",
	Short: "Skemr CLI is a command-line tool for Skemr",
	Long:  `Skemr CLI allows you to interact with the Skemr platform from the command line.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Skemr CLI! Use --help to see available commands.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(os.Stderr, err)
		os.Exit(1)
	}
}
