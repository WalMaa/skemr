package rulengn

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"slices"
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

func (r *RuleEngine) ProcessMigrationFiles(c context.Context, statements []MigrationFileDto, rules []models.Rule, entities []models.DatabaseEntity) ([]StatementResult, error) {
	results := make(chan StatementResult, len(statements))
	var wg sync.WaitGroup
	for _, statement := range statements {
		stmt := statement
		wg.Add(1)
		go func() {
			defer wg.Done()
			slog.Debug("Processing migration file", "file", stmt.File)
			stmtResults, err := r.CheckStatement(stmt, rules, entities)
			if err != nil {
				slog.Error("Error checking statement", slog.String("statement", stmt.File), slog.String("error", err.Error()))
				return
			}

			for _, res := range stmtResults {
				select {
				case results <- res:
					slog.Debug("Statement result sent", slog.String("statement", stmt.File), slog.String("rule", res.Rule.Name))
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
func (r *RuleEngine) CheckStatement(migrationFileDto MigrationFileDto, rules []models.Rule, entities []models.DatabaseEntity) ([]StatementResult, error) {
	slog.Debug("Checking migration file", "file", migrationFileDto.File)

	file, err := os.ReadFile(migrationFileDto.File)
	if err != nil {
		slog.Error("Error reading migration file", "file", migrationFileDto.File, "err", err)
		return nil, err
	}

	statementActions, err := parser.ParseSql(string(file))
	statementResults := make([]StatementResult, 0)

	if err != nil {
		slog.Error("Error parsing migrationFileDto", "migrationFileDto", migrationFileDto, "err", err)
		return nil, err
	}

	for _, rule := range rules {
		slog.Debug("Evaluating rule against migration file", slog.String("rule_name", rule.Name), slog.String("migration_file", migrationFileDto.File))
		for _, action := range statementActions {
			// If the database entity is a column, it is not enough to match the name of the column in the rule with the name of the column in the statement action,
			// we also need to check if the columns are in the same table. This is because there could be multiple columns with the same name in different tables,
			// and we don't want to trigger a rule violation if the column in the statement action is not the same as the column in the rule.
			if rule.DataBaseEntity.Type == models.DatabaseEntityTypeColumn {
				// Get the parent database entity (table) for the column in the rule
				i := slices.IndexFunc(entities, func(entity models.DatabaseEntity) bool {
					return entity.ID == *rule.DataBaseEntity.ParentId
				})

				if i == -1 {
					slog.Warn("Parent database entity not found for rule", slog.String("rule_name", rule.Name), slog.String("parent_id", rule.DataBaseEntity.ParentId.String()))
					continue
				}

				parentEntity := entities[i]
				slog.Debug("Found parent database entity for rule", slog.String("rule_name", rule.Name), slog.String("parent_entity_name", parentEntity.Name))

				// Check if the parent entity name matches the table name in the statement action
				if parentEntity.Name != action.Relation {
					slog.Debug("Parent entity name does not match action table, skipping rule evaluation for this action", slog.String("rule_name", rule.Name), slog.String("parent_entity_name", parentEntity.Name), slog.String("action_table", action.Relation))
					continue
				}
			}
			if rule.DataBaseEntity.Name == action.Target {
				slog.Debug("Rule target matches migrationFileDto target", slog.String("rule_database_entity", rule.DataBaseEntity.Name), slog.String("statement_target", action.Target))
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
	return StatementResult{
		Type:  models.RuleTypeLocked,
		Rule:  rule,
		File:  statementDto.File,
		Error: err,
	}
}

func (r *RuleEngine) warnAction(rule models.Rule, statementDto MigrationFileDto) StatementResult {
	return StatementResult{
		Type: models.RuleTypeWarn,
		Rule: rule,
		File: statementDto.File,
	}
}

func (r *RuleEngine) advisoryAction(rule models.Rule, statementDto MigrationFileDto) StatementResult {
	return StatementResult{
		Type: models.RuleTypeAdvisory,
		Rule: rule,
		File: statementDto.File,
	}
}

func (r *RuleEngine) deprecatedAction(rule models.Rule, statementDto MigrationFileDto) StatementResult {
	return StatementResult{
		Type: models.RuleTypeDeprecated,
		Rule: rule,
		File: statementDto.File,
	}
}
