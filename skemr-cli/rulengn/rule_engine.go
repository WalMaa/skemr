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

func (r *RuleEngine) ProcessStatements(c context.Context, statements []string, database *models.Database) (chan bool, error) {
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
func (r *RuleEngine) CheckStatement(c context.Context, statement string, database *models.Database) bool {
	slog.Info("CheckStatement", slog.String("statement", statement))
	_, err := parser.ParseSql(statement)

	if err != nil {
		slog.Error("Error parsing statement", "statement", statement, "err", err)
		return false
	}

	// TODO: Check against rules
	//args := sqlc.ListRulesByCriteriaParams{
	//	DatabaseID:   database.ID,
	//	Scope:        sqlc.RuleScopeTable,
	//	RelationName: pgtype.Text{String: stmtact.Relation, Valid: true},
	//	Target:       stmtact.Target,
	//}
	//rules, err := r.db.ListRulesByCriteria(c, args)

	return true
}
