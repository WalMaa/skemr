package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/walmaa/skemr-cli/controlplaneclient"
	"github.com/walmaa/skemr-cli/rulengn"
)

func init() {
	validateCmd.Flags().StringP("projectId", "P", "", "ID of your project")
	err := validateCmd.MarkFlagRequired("projectId")
	if err != nil {
		slog.Error("Unable to mark projectId as required")
	}
	validateCmd.Flags().StringP("databaseId", "D", "", "ID of your database")
	rootCmd.AddCommand(validateCmd)
	err = validateCmd.MarkFlagRequired("databaseId")
	if err != nil {
		slog.Error("Unable to mark databaseId as required")
	}

}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate SQL statements",
	Long:  `Validate SQL statements for correctness and compliance with Skemr rules.`,
	Run: func(cmd *cobra.Command, args []string) {
		projectId, err := cmd.Flags().GetString("projectId")
		if err != nil {
			slog.Error("projectId not provided")
			os.Exit(1)
		}
		databaseId, err := cmd.Flags().GetString("databaseId")
		if err != nil {
			slog.Error("databaseId not provided")
			os.Exit(1)
		}
		slog.Info("Validate command executed", slog.String("project_id", projectId))
		cmd.Context()
		ruleEngine := rulengn.NewRuleEngine()

		// Get rules
		rules := controlplaneclient.GetRules(cmd.Context(), projectId, databaseId)

		// Process files
		filePaths := make([]string, 0)
		collectFilePathsFromDir(&filePaths, "./skemr-cli/test")

		// Rule check
		dtos := make([]rulengn.MigrationFileDto, len(filePaths))
		for i, path := range filePaths {
			dtos[i] = rulengn.MigrationFileDto{File: path}
		}
		_, err = ruleEngine.ProcessMigrationFiles(cmd.Context(), dtos, rules)
		if err != nil {
			err = fmt.Errorf("Error while validating migrations")
			fmt.Println(err)
			os.Exit(1)
		}
		slog.Info("Processed migration files", "fileCount", len(filePaths))
	},
}

func collectFilePathsFromDir(filePaths *[]string, dirPath string) {
	slog.Debug("Gathering filepaths in directory", slog.String("directory", dirPath))
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
		// Recursively add files from subdirectories
		if !entry.IsDir() {
			slog.Debug("Adding file", slog.String("file", entry.Name()))
			*filePaths = append(*filePaths, filepath.Join(path, entry.Name()))
		} else {
			collectFilePathsFromDir(filePaths, filepath.Join(dirPath, entry.Name()))
		}
	}
}
