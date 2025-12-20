package rulengn

import (
	"context"
	"log/slog"
	"sync"

	"github.com/walmaa/skemr-cli/parser"
	"github.com/walmaa/skemr-common/models"
)

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
			results <- r.CheckStatement(c, stmt, rules)
		}()

	}

	go func() {
		wg.Wait()
		close(results)
	}()
	return results, nil
}

// CheckStatement checks if the given SQL statement matches any rules in the database for the specified project.
func (r *RuleEngine) CheckStatement(c context.Context, statement string, rules []models.Rule) bool {
	slog.Info("CheckStatement", slog.String("statement", statement))
	statementAction, err := parser.ParseSql(statement)

	if err != nil {
		slog.Error("Error parsing statement", "statement", statement, "err", err)
		return false
	}

	for _, rule := range rules {
		if rule.Target == statementAction.Target {
			slog.Info("Rule target matches statement target", slog.String("rule_target", rule.Target), slog.String("statement_target", statementAction.Target))

		}
	}

	return true
}
