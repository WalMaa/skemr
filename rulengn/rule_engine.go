package rulengn

import (
	"context"
	"fmt"
	"github.com/walmaa/skemr/db/sqlc"
	"github.com/walmaa/skemr/parser"
	"log/slog"
	"sync"
)

type RuleEngine struct {
	db sqlc.Querier
}

func NewRuleEngine(q sqlc.Querier) *RuleEngine {
	return &RuleEngine{db: q}
}

func (r *RuleEngine) ProcessStatements(c context.Context, statements []string, database *sqlc.Database) (chan bool, error) {
	results := make(chan bool, len(statements))
	var wg sync.WaitGroup
	for _, statement := range statements {
		stmt := statement
		wg.Add(1)
		go func() {
			defer wg.Done()
			slog.Info("Processing statement", slog.String("statement", stmt))
			results <- r.CheckStatement(c, stmt, database)
		}()

	}

	go func() {
		wg.Wait()
		close(results)
	}()
	return results, nil
}

// CheckStatement checks if the given SQL statement matches any rules in the database for the specified project.
func (r *RuleEngine) CheckStatement(c context.Context, statement string, database *sqlc.Database) bool {
	slog.Info("CheckStatement", slog.String("statement", statement))
	stmtact, err := parser.ParseSql(statement)

	if err != nil {
		slog.Error("Error parsing statement", "statement", statement, "err", err)
		return false
	}

	args := sqlc.ListRulesByCriteriaParams{
		DatabaseID:   database.ID,
		Scope:        sqlc.RuleScopeTable,
		RelationName: &stmtact.Relation,
		Target:       stmtact.Target,
	}
	rules, err := r.db.ListRulesByCriteria(c, args)

	fmt.Println(rules)
	return true
}
