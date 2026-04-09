package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/walmaa/skemr-cli/logger"
)

func init() {
	rootCmd.PersistentFlags().String("log-level", "info", "Set the logging level (debug, info, warn, error)")
	_ = viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))

	viper.SetDefault("log-level", "warn")
}

var rootCmd = &cobra.Command{
	Use:   "skemr-cli",
	Short: "Skemr CLI is a command-line tool for Skemr",
	Long:  `Skemr CLI allows you to interact with the Skemr platform from the command line.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var level slog.Level
		if err := level.UnmarshalText([]byte(viper.GetString("log-level"))); err != nil {
			fmt.Printf("Invalid log level: %s\n", viper.GetString("log-level"))
			os.Exit(1)
		}
		logger.BuildLogger(level)
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Skemr CLI! Use --help to see available commands.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
