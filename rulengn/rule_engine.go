package rulengn

import (
	"context"
	"log/slog"
	"skemr/db/sqlc"
	"skemr/parser"
)

type RuleEngine struct {
	db sqlc.Querier
}

func (r *RuleEngine) CheckStatement(c context.Context, statement string, project *sqlc.Project) bool {
	stmtact, err := parser.ParseSql(statement)

	if err != nil {
		slog.Error("Error parsing statement", "statement", statement, "err", err)
		return false
	}

	args := sqlc.ListRulesByCriteriaParams{
		ProjectID: project.ID,
		Scope:     sqlc.,
	}
	r.db.ListRulesByCriteria(c, args)

}
