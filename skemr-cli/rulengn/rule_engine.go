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

type SqlStatementDto struct {
	Statement string
	File      string
	Line      int
}

type RuleEngine struct {
}

func NewRuleEngine() *RuleEngine {
	return &RuleEngine{}
}

func (r *RuleEngine) ProcessStatements(c context.Context, statements []SqlStatementDto, rules []models.Rule) ([]StatementResult, error) {
	results := make(chan StatementResult, len(statements))
	var wg sync.WaitGroup
	for _, statement := range statements {
		stmt := statement
		wg.Add(1)
		go func() {
			defer wg.Done()
			slog.Info("Processing statement", slog.String("statement", stmt.Statement))
			stmtResult, err := r.CheckStatement(stmt, rules)
			if err != nil {
				slog.Error("Error checking statement", slog.String("statement", stmt.Statement), slog.String("error", err.Error()))
				return
			}

			for _, res := range stmtResult {
				select {
				case results <- res:
					slog.Info("Statement result sent", slog.String("statement", stmt.Statement), slog.String("rule", res.Rule.Name))
				case <-c.Done():
					slog.Warn("Context done before sending result", slog.String("statement", stmt.Statement))
					return
				}
			}
		}()

	}

	go func() {
		wg.Wait()
		close(results)

	}()
	resultsSlice := make([]StatementResult, 0)
	for res := range results {
		resultsSlice = append(resultsSlice, res)
	}

	return resultsSlice, nil
}

// CheckStatement checks if the given SQL statement matches any rules in the database for the specified project.
func (r *RuleEngine) CheckStatement(statementDto SqlStatementDto, rules []models.Rule) ([]StatementResult, error) {
	slog.Info("CheckStatement", slog.String("statementDto", statementDto.Statement))
	statementAction, err := parser.ParseSql(statementDto.Statement)
	statementResults := make([]StatementResult, 0)

	if err != nil {
		slog.Error("Error parsing statementDto", "statementDto", statementDto, "err", err)
		return nil, err
	}

	for _, rule := range rules {
		if rule.Target == statementAction.Target {
			slog.Info("Rule target matches statementDto target", slog.String("rule_target", rule.Target), slog.String("statement_target", statementAction.Target))
			switch rule.Type {
			case models.RuleTypeLocked:
				violation := r.lockAction(rule, statementAction, statementDto)
				statementResults = append(statementResults, violation)
			case models.RuleTypeWarn:
				warning := r.warnAction(rule, statementDto)
				statementResults = append(statementResults, warning)
			case models.RuleTypeAdvisory:
				advisory := r.advisoryAction(rule, statementDto)
				statementResults = append(statementResults, advisory)
			case models.RuleTypeDeprecated:
				warning := r.deprecatedAction(rule, statementDto)
				statementResults = append(statementResults, warning)
			default:
				slog.Warn("Unknown rule type", slog.String("rule_type", string(rule.Type)))

			}
		}
	}
	return statementResults, nil
}

// lockAction handles rule matches where the rule type was defined as "locked"
func (r *RuleEngine) lockAction(rule models.Rule, statementAction parser.StatementAction, statementDto SqlStatementDto) StatementResult {
	err := fmt.Errorf("Lock rule violated: rule %q violated by statementDto %q. %q is not allowed on %q ", rule.Name, statementAction, statementAction.Action, statementAction.Target)
	fmt.Fprintln(os.Stderr, err)
	return StatementResult{
		Type:      models.RuleTypeLocked,
		Rule:      rule,
		Statement: statementDto.Statement,
		Error:     err,
		File:      statementDto.File,
		Line:      statementDto.Line,
	}
}

func (r *RuleEngine) warnAction(rule models.Rule, statementDto SqlStatementDto) StatementResult {
	fmt.Println("Warning: ", rule.Name, " triggered by statement: ", statementDto.Statement)
	return StatementResult{
		Type:      models.RuleTypeWarn,
		Rule:      rule,
		Statement: statementDto.Statement,
		File:      statementDto.File,
		Line:      statementDto.Line,
	}
}

func (r *RuleEngine) advisoryAction(rule models.Rule, statementDto SqlStatementDto) StatementResult {
	fmt.Println("Advisory: ", rule.Name, " triggered by statement: ", statementDto.Statement)
	return StatementResult{
		Type:      models.RuleTypeAdvisory,
		Rule:      rule,
		Statement: statementDto.Statement,
		File:      statementDto.File,
		Line:      statementDto.Line,
	}
}

func (r *RuleEngine) deprecatedAction(rule models.Rule, statementDto SqlStatementDto) StatementResult {
	fmt.Println("Deprecated: ", rule.Name, " triggered by statement: ", statementDto.Statement)
	return StatementResult{
		Type:      models.RuleTypeDeprecated,
		Rule:      rule,
		Statement: statementDto.Statement,
		File:      statementDto.File,
		Line:      statementDto.Line,
	}
}
