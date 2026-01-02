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

type Violation struct {
	Rule      models.Rule
	Statement string
	File      string
	Line      int
	Error     error
}

type RuleEngine struct {
}

func NewRuleEngine() *RuleEngine {
	return &RuleEngine{}
}

func (r *RuleEngine) ProcessStatements(c context.Context, statements []string, rules []models.Rule) (chan bool, error) {
	results := make(chan bool, len(statements))
	var wg sync.WaitGroup
	for _, statement := range statements {
		stmt := statement
		wg.Add(1)
		go func() {
			defer wg.Done()
			slog.Info("Processing statement", slog.String("statement", stmt))
			results <- r.CheckStatement(stmt, rules)
		}()

	}

	go func() {
		wg.Wait()
		close(results)
	}()
	return results, nil
}

// CheckStatement checks if the given SQL statement matches any rules in the database for the specified project.
func (r *RuleEngine) CheckStatement(statement string, rules []models.Rule) (chan error, error) {
	slog.Info("CheckStatement", slog.String("statement", statement))
	statementAction, err := parser.ParseSql(statement)
	violations := make(chan error, 1)

	defer close(violations)

	if err != nil {
		slog.Error("Error parsing statement", "statement", statement, "err", err)
		return nil, err
	}

	for _, rule := range rules {
		if rule.Target == statementAction.Target {
			slog.Info("Rule target matches statement target", slog.String("rule_target", rule.Target), slog.String("statement_target", statementAction.Target))
			if rule.Type == models.RuleTypeLocked {
				violations <- r.lockAction(rule, statementAction)
			}
		}
	}
	return violations, nil
}

// lockAction handles rule matches where the rule type was defined as "locked"
func (r *RuleEngine) lockAction(rule models.Rule, statementAction parser.StatementAction) error {
	err := fmt.Errorf("Lock rule violated: rule %q violated by statement %q. %q is not allowed on %q ", rule.Name, statementAction, statementAction.Action, statementAction.Target)
	fmt.Fprintln(os.Stderr, err)
	return err
}
