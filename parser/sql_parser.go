package parser

import (
	"fmt"
	pg_query "github.com/pganalyze/pg_query_go/v6"
	"log"
	skemr "skemr/db"
)

type StatementAction struct {
	Target   string
	Action   SqlAction
	Relation string
}

type SqlAction string

const (
	SqlActionDropColumn   SqlAction = "DROP COLUMN"
	SqlActionRenameColumn SqlAction = "RENAME COLUMN"
	SqlActionUndefined    SqlAction = "UNDEFINED"
)

/*
ParseSql parses the SQL statement and returns a structured representation of the SQL.
*/
func ParseSql(sql string) (StatementAction, error) {
	tree, err := pg_query.Parse(sql)
	if err != nil {
		log.Fatal(err)
		return StatementAction{}, err
	}
	stmts := tree.Stmts

	if len(stmts) == 0 {
		return StatementAction{}, nil
	}

	node := stmts[0].GetStmt()

	if node.GetAlterTableStmt() != nil {
		return parseAlterTable(node)
	}

	if node.GetRenameStmt() != nil {
		return parseRenameStmt(node)
	}

	// If the statement is not recognized, return an error
	return StatementAction{
		Target:   "",
		Action:   SqlActionUndefined,
		Relation: "",
	}, fmt.Errorf("unsupported SQL statement: %s", sql)
}

func parseRenameStmt(node *pg_query.Node) (StatementAction, error) {
	renameStmt := node.GetRenameStmt()
	relName := renameStmt.Relation.Relname
	target := renameStmt.Subname
	action := SqlActionRenameColumn

	return StatementAction{
		Target:   target,
		Action:   action,
		Relation: relName,
	}, nil
}

func parseAlterTable(node *pg_query.Node) (StatementAction, error) {
	alterTable := node.GetAlterTableStmt()
	relName := alterTable.Relation.Relname
	target := ""
	action := SqlActionUndefined

	for _, cmd := range alterTable.Cmds {
		target = cmd.GetAlterTableCmd().Name
		if cmd.GetAlterTableCmd().GetSubtype() == pg_query.AlterTableType_AT_DropColumn {
			action = SqlActionDropColumn
		}
	}

	return StatementAction{
		Target:   target,
		Action:   action,
		Relation: relName,
	}, nil
}

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
	alterTable := node.GetAlterTableStmt()
	log.Println(alterTable)
	cmd := node.GetAlterTableCmd()
	relname := alterTable.Relation.Relname
	cmds := alterTable.Cmds
	for _, c := range cmds {
		log.Println(c)
		log.Println(c.GetAlterTableCmd().Name)
		log.Println(c.GetAlterTableCmd().GetSubtype())
	}

	log.Println(relname)
	log.Println(cmd)
	return true
}

func parseStatement(node *pg_query.Node) {
	x := node.GetAExpr()
	log.Println(x)

}
