package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/walmaa/skemr-cli/controlplaneclient"
	"github.com/walmaa/skemr-cli/reporter"
	"github.com/walmaa/skemr-cli/rulengn"
)

var (
	projectId         string
	databaseId        string
	token             string
	migrationFilesDir string
)

func init() {
	validateCmd.Flags().StringVarP(&projectId, "projectId", "P", "", "ID of your project")
	validateCmd.Flags().StringVarP(&databaseId, "databaseId", "D", "", "ID of your database")
	validateCmd.Flags().StringVarP(&token, "token", "T", "", "API Token")
	validateCmd.Flags().StringVarP(&migrationFilesDir, "migrationFilesDir", "M", "migrations", "Directory containing migration files to validate")
	validateCmd.Flags().String("host", "https://api.skemr.com", "URL of the Skemr control plane")
	viper.BindPFlag(
		"controlPlaneUrl",
		validateCmd.Flags().Lookup("host"),
	)

	mustMarkRequired(validateCmd, "projectId", "databaseId", "token", "migrationFilesDir")

	rootCmd.AddCommand(validateCmd)
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate SQL statements",
	Long:  `Validate SQL statements for correctness and compliance with Skemr rules.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := cmd.Context()

		slog.Debug("Validate command executed", "project_id", projectId, "database_id", databaseId, "migration_files_dir", migrationFilesDir)
		ruleEngine := rulengn.NewRuleEngine()

		// Get rules
		rules, err := controlplaneclient.GetRules(c, projectId, databaseId, token)

		if err != nil {
			slog.Error("Error fetching rules from control plane", "err", err)
			os.Exit(1)
		}
		slog.Debug("Fetched rules from control plane", "ruleCount", len(rules))

		// Fetch all database entities. This is used to match columns to right parents (tables)
		entities, err := controlplaneclient.GetDatabaseEntities(c, projectId, databaseId, token)
		slog.Debug("Fetched database entities from control plane", "entityCount", len(entities))

		if err != nil {
			os.Exit(1)
		}

		// Process files
		filePaths := make([]string, 0)
		collectFilePathsFromDir(&filePaths, migrationFilesDir)

		// Rule check
		dtos := make([]rulengn.MigrationFileDto, len(filePaths))
		for i, path := range filePaths {
			dtos[i] = rulengn.MigrationFileDto{File: path}
		}
		statementResults, err := ruleEngine.ProcessMigrationFiles(c, dtos, rules, entities)
		if err != nil {
			err = fmt.Errorf("error while validating migrations")
			slog.Error("Error while validating migrations", "err", err)
			os.Exit(1)
		}

		reporter.PrintSummary(statementResults)

		// Check if any results are of type error and exit with non-zero code if so
		hasErrors := false
		for _, res := range statementResults {
			if res.Error != nil {
				hasErrors = true
				slog.Debug("Validation error encountered", "file", res.File, "line", res.Line, "error", res.Error)
			} else {
				slog.Debug("Validation result", "file", res.File, "line", res.Line, "rule", res.Rule.Name, "type", res.Type)
			}
		}

		if hasErrors {
			os.Exit(1)
		}

	},
}

func collectFilePathsFromDir(filePaths *[]string, dirPath string) {
	slog.Debug("Gathering filepaths in directory", slog.String("directory", dirPath))
	cwd, err := os.Getwd()
	fileExtensions := []string{".sql"}

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
			ext := filepath.Ext(entry.Name())
			if !slices.Contains(fileExtensions, ext) {
				slog.Debug("Skipping file with unsupported extension", slog.String("file", entry.Name()), slog.String("extension", ext))
				continue
			}
			slog.Debug("Adding file", slog.String("file", entry.Name()))
			*filePaths = append(*filePaths, filepath.Join(path, entry.Name()))
		} else {
			collectFilePathsFromDir(filePaths, filepath.Join(dirPath, entry.Name()))
		}
	}
}

func mustMarkRequired(cmd *cobra.Command, flags ...string) {
	for _, f := range flags {
		if err := cmd.MarkFlagRequired(f); err != nil {
			slog.Error("Unable to mark flag as required", "flag", f, "err", err)
		}
	}
}
