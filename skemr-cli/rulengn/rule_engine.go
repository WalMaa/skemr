package rulengn

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/walmaa/skemr-cli/parser"
	"github.com/walmaa/skemr-common/models"
)

type StatementResult struct {
	Type      models.RuleType
	Statement string
	File      string
	Line      int
	Rule      models.Rule
	Error     error
}

type MigrationFileDto struct {
	File string
}

type RuleEngine struct {
}

func NewRuleEngine() *RuleEngine {
	return &RuleEngine{}
}

func (r *RuleEngine) ProcessMigrationFiles(c context.Context, statements []MigrationFileDto, rules []models.Rule) ([]StatementResult, error) {
	results := make(chan StatementResult, len(statements))
	var wg sync.WaitGroup
	for _, statement := range statements {
		stmt := statement
		wg.Add(1)
		go func() {
			defer wg.Done()
			slog.Info("Processing migration file", "file", stmt.File)
			stmtResult, err := r.CheckStatement(stmt, rules)
			if err != nil {
				slog.Error("Error checking statement", slog.String("statement", stmt.File), slog.String("error", err.Error()))
				return
			}

			for _, res := range stmtResult {
				select {
				case results <- res:
					slog.Info("Statement result sent", slog.String("statement", stmt.File), slog.String("rule", res.Rule.Name))
				case <-c.Done():
					slog.Warn("Context done before sending result", slog.String("statement", stmt.File))
					return
				}
			}
		}()

	}

	go func() {
		wg.Wait()
		close(results)

	}()
	// Collect chan to slice
	resultsSlice := make([]StatementResult, 0)
	for res := range results {
		resultsSlice = append(resultsSlice, res)
	}

	return resultsSlice, nil
}

// CheckStatement checks if the given SQL statement matches any rules in the database for the specified project.
func (r *RuleEngine) CheckStatement(migrationFileDto MigrationFileDto, rules []models.Rule) ([]StatementResult, error) {
	slog.Info("Checking migration file", "file", migrationFileDto.File)

	file, err := os.ReadFile(migrationFileDto.File)
	if err != nil {
		return nil, err
	}

	statementActions, err := parser.ParseSql(string(file))
	statementResults := make([]StatementResult, 0)

	if err != nil {
		slog.Error("Error parsing migrationFileDto", "migrationFileDto", migrationFileDto, "err", err)
		return nil, err
	}

	for _, rule := range rules {
		for _, action := range statementActions {

			if rule.DataBaseEntity.Name == action.Target {
				slog.Info("Rule target matches migrationFileDto target", slog.String("rule_database_entity", rule.DataBaseEntity.Name), slog.String("statement_target", action.Target))
				switch rule.RuleType {
				case models.RuleTypeLocked:
					violation := r.lockAction(rule, action, migrationFileDto)
					statementResults = append(statementResults, violation)
				case models.RuleTypeWarn:
					warning := r.warnAction(rule, migrationFileDto)
					statementResults = append(statementResults, warning)
				case models.RuleTypeAdvisory:
					advisory := r.advisoryAction(rule, migrationFileDto)
					statementResults = append(statementResults, advisory)
				case models.RuleTypeDeprecated:
					warning := r.deprecatedAction(rule, migrationFileDto)
					statementResults = append(statementResults, warning)
				default:
					slog.Warn("Unknown rule type", slog.String("rule_type", string(rule.RuleType)))

				}
			}
		}

	}
	return statementResults, nil
}

// lockAction handles rule matches where the rule type was defined as "locked"
func (r *RuleEngine) lockAction(rule models.Rule, statementAction parser.StatementAction, statementDto MigrationFileDto) StatementResult {
	err := fmt.Errorf("Lock rule violated: rule %q violated by action %q on target %q", rule.Name, statementAction.Action, statementAction.Target)
	fmt.Fprintln(os.Stderr, err)
	return StatementResult{
		Type:  models.RuleTypeLocked,
		Rule:  rule,
		File:  statementDto.File,
		Error: err,
	}
}

func (r *RuleEngine) warnAction(rule models.Rule, statementDto MigrationFileDto) StatementResult {
	fmt.Println("Warning:", rule.Name, "triggered by file:", statementDto.File)
	return StatementResult{
		Type: models.RuleTypeWarn,
		Rule: rule,
		File: statementDto.File,
	}
}

func (r *RuleEngine) advisoryAction(rule models.Rule, statementDto MigrationFileDto) StatementResult {
	fmt.Println("Advisory:", rule.Name, "triggered by file:", statementDto.File)
	return StatementResult{
		Type: models.RuleTypeAdvisory,
		Rule: rule,
		File: statementDto.File,
	}
}

func (r *RuleEngine) deprecatedAction(rule models.Rule, statementDto MigrationFileDto) StatementResult {
	fmt.Println("Deprecated:", rule.Name, "triggered by file:", statementDto.File)
	return StatementResult{
		Type: models.RuleTypeDeprecated,
		Rule: rule,
		File: statementDto.File,
	}
}
