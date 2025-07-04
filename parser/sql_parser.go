package parser

import (
	pg_query "github.com/pganalyze/pg_query_go/v6"
	"log"
	skemr "skemr/db"
)

func ParseRule(rule *skemr.Rule, sql string) (string, error) {
	log.Printf("Rule: %#v", rule)
	tree, err := pg_query.Parse(sql)

	log.Println(tree)
	stmts := tree.Stmts
	log.Println(stmts)

	for _, stmt := range stmts {
		checkRule(rule, stmt.GetStmt())
	}

	if err != nil {
		log.Fatal(err)
	}

	return sql, nil
}

func checkRule(rule *skemr.Rule, node *pg_query.Node) bool {
	selectCase := node.GetAlterTableStmt()
	log.Println(selectCase)
	x := node.GetAExpr()
	log.Println(x)

	return true
}

func parseStatement(node *pg_query.Node) {
	x := node.GetAExpr()
	log.Println(x)

}
