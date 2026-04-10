package reporter

import (
	"fmt"
	"os"

	"charm.land/lipgloss/v2"
	"github.com/walmaa/skemr-cli/rulengn"
	"github.com/walmaa/skemr-common/models"
)

type ValidationSummary struct {
	DeprecatedCount int
	WarningCount    int
	AdvisoryCount   int
	ErrorCount      int
}

var (
	errorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
	warningStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true)
	advisoryStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	deprecatedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true)
)

func severityLabel(ruleType models.RuleType, colorize bool) string {
	var (
		label string
		style lipgloss.Style
	)

	switch ruleType {
	case models.RuleTypeLocked:
		label = "Error"
		style = errorStyle
	case models.RuleTypeWarn:
		label = "Warning"
		style = warningStyle
	case models.RuleTypeAdvisory:
		label = "Advisory"
		style = advisoryStyle
	case models.RuleTypeDeprecated:
		label = "Deprecated"
		style = deprecatedStyle
	default:
		return "Unknown"
	}

	if !colorize {
		return label
	}

	return style.Render(label)
}

// PrintSummary creates a summary of the validation results, counting how many statements triggered each rule type and how many errors occurred, and logs it.
func PrintSummary(statementResults []rulengn.StatementResult, rules []models.Rule) {

	fmt.Fprintf(os.Stdout, "Number of rules checked: %d\n", len(rules))

	fmt.Fprintf(os.Stdout, "\nValidation Summary:\n")

	summary := ValidationSummary{}

	for _, res := range statementResults {
		label := severityLabel(res.Type, true)
		switch res.Type {
		case models.RuleTypeDeprecated:
			summary.DeprecatedCount++
			fmt.Fprintf(os.Stdout, "%s: %s (File: %s)\n", label, res.Rule.Name, res.File)
		case models.RuleTypeWarn:
			summary.WarningCount++
			fmt.Fprintf(os.Stdout, "%s: %s (File: %s)\n", label, res.Rule.Name, res.File)
		case models.RuleTypeLocked:
			summary.ErrorCount++
			fmt.Fprintf(os.Stdout, "%s: rule \"%s\" violated (File: %s)\n", label, res.Rule.Name, res.File)
		case models.RuleTypeAdvisory:
			summary.AdvisoryCount++
			fmt.Fprintf(os.Stdout, "%s: %s (File: %s)\n", label, res.Rule.Name, res.File)
		}
	}

	fmt.Fprintf(os.Stdout, "Deprecated: %d\n", summary.DeprecatedCount)
	fmt.Fprintf(os.Stdout, "Warnings: %d\n", summary.WarningCount)
	fmt.Fprintf(os.Stdout, "Advisories: %d\n", summary.AdvisoryCount)
	fmt.Fprintf(os.Stdout, "Errors: %d\n", summary.ErrorCount)

}
