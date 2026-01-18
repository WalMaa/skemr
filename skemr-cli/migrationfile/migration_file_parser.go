package migrationfile

import (
	"log/slog"
	"os"

	"github.com/walmaa/skemr-cli/rulengn"
)

// Builds a SqlStatementDro
func GetSqlStatements(filePath string) rulengn.SqlStatementDto {
	// Get file
	f, err := os.ReadFile(filePath)
	if err != nil {
		slog.Error("Error reading file", "filePath", filePath)
		panic(err)
	}

	// file as string
	str := string(f)
}
