package cmd

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/walmaa/skemr-cli/rulengn"
)

func init() {
	rootCmd.AddCommand(validateCmd)
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate SQL statements",
	Long:  `Validate SQL statements for correctness and compliance with Skemr rules.`,
	Run: func(cmd *cobra.Command, args []string) {
		projectID := ""
		slog.Info("Validate command executed", slog.String("project_id", projectID))
		cmd.Context()
		ruleEngine := rulengn.NewRuleEngine()
		ruleEngine.ProcessStatements(cmd.Context(), nil, nil)
		printFilesInDirectory("")
	},
}

func printFilesInDirectory(dirPath string) {
	slog.Info("Printing files in directory", slog.String("directory", dirPath))
	cwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	path := filepath.Join(cwd, dirPath)
	dat, err := os.ReadDir(path)

	if err != nil {
		panic(err)
	}

	for _, entry := range dat {
		slog.Info(entry.Name())
	}
}
