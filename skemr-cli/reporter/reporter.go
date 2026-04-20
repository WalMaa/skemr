package reporter

import (
	"fmt"
	"os"
	"strings"

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
	const dividerWidth = 72
	fmt.Fprintf(os.Stdout, "%s\n", strings.Repeat("-", dividerWidth))
	fmt.Fprintln(os.Stdout, "Validation Summary")
	fmt.Fprintf(os.Stdout, "%s\n", strings.Repeat("-", dividerWidth))

	summary := ValidationSummary{}
	warnings := make([]rulengn.StatementResult, 0)
	advisories := make([]rulengn.StatementResult, 0)
	errors := make([]rulengn.StatementResult, 0)
	deprecated := make([]rulengn.StatementResult, 0)

	for _, res := range statementResults {
		switch res.Type {
		case models.RuleTypeDeprecated:
			summary.DeprecatedCount++
			deprecated = append(deprecated, res)
		case models.RuleTypeWarn:
			summary.WarningCount++
			warnings = append(warnings, res)
		case models.RuleTypeLocked:
			summary.ErrorCount++
			errors = append(errors, res)
		case models.RuleTypeAdvisory:
			summary.AdvisoryCount++
			advisories = append(advisories, res)
		}
	}

	fmt.Fprintf(os.Stdout, "%-16s %d\n", "Rules checked: ", len(rules))
	fmt.Fprintf(os.Stdout, "%-16s %d\n\n", "Findings: ", len(statementResults))

	printSection("Warnings", warnings, false)
	fmt.Fprintln(os.Stdout)
	printSection("Advisories", advisories, false)
	fmt.Fprintln(os.Stdout)
	printSection("Errors", errors, true)

	if len(deprecated) > 0 {
		fmt.Fprintln(os.Stdout)
		printSection("Deprecated", deprecated, false)
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "Totals")
	fmt.Fprintf(os.Stdout, "  %-14s %d\n", "Deprecated:", summary.DeprecatedCount)
	fmt.Fprintf(os.Stdout, "  %-14s %d\n", "Warnings:", summary.WarningCount)
	fmt.Fprintf(os.Stdout, "  %-14s %d\n", "Advisories:", summary.AdvisoryCount)
	fmt.Fprintf(os.Stdout, "  %-14s %d\n", "Errors:", summary.ErrorCount)

}

func printSection(title string, results []rulengn.StatementResult, isError bool) {
	fmt.Fprintf(os.Stdout, "%-15s (%d)\n", title, len(results))
	for _, res := range results {
		message := res.Rule.Name
		if isError {
			message = fmt.Sprintf("rule \"%s\" violated", res.Rule.Name)
		}
		fmt.Fprintf(os.Stdout, "  -  %-30s by statement \"%s\" in file: %s\n", message, res.Statement, res.File)
	}

}
