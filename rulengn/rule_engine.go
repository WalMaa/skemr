package rulengn

import (
	"context"
	"fmt"
	"log/slog"
	"skemr/db/sqlc"
	"skemr/parser"
)

type RuleEngine struct {
	db sqlc.Querier
}

func NewRuleEngine(q sqlc.Querier) *RuleEngine {
	return &RuleEngine{db: q}
}

// CheckStatement checks if the given SQL statement matches any rules in the database for the specified project.
func (r *RuleEngine) CheckStatement(c context.Context, statement string, project *sqlc.Project) bool {
	stmtact, err := parser.ParseSql(statement)

	if err != nil {
		slog.Error("Error parsing statement", "statement", statement, "err", err)
		return false
	}

	args := sqlc.ListRulesByCriteriaParams{
		ProjectID:    project.ID,
		Scope:        sqlc.RuleScopeTable,
		RelationName: &stmtact.Relation,
		Target:       stmtact.Target,
	}
	rules, err := r.db.ListRulesByCriteria(c, args)

	fmt.Println(rules)
	return true
}
