package test

import (
	skemr "skemr/db"
	"skemr/parser"
	"testing"
)

func TestSQLParser(t *testing.T) {
	rule := &skemr.Rule{
		ID:     1,
		Name:   "Test Rule",
		Type:   skemr.RuleTypeLock,
		Scope:  skemr.RuleScopeColumn,
		Target: "name",
	}
	sql := "ALTER TABLE rules DROP COLUMN name"
	_, _ = parser.ParseRule(rule, sql)

}

func TestParseSqlDropColumn(t *testing.T) {
	sql := "ALTER TABLE rules DROP COLUMN name"
	statementAction, err := parser.ParseSql(sql)

	if err != nil {
		t.Fatalf("Failed to parse SQL: %v", err)
	}
	if statementAction.Target != "name" || statementAction.Action != parser.SqlActionDropColumn || statementAction.Relation != "rules" {
		t.Fatalf("Expected target 'name', action 'DROP', relation 'rules', got target '%s', action '%s', relation '%s'", statementAction.Target, statementAction.Action, statementAction.Relation)
	}
}

func TestParseSqlRenameColumn(t *testing.T) {
	sql := "ALTER TABLE rules RENAME COLUMN name TO new_name"
	statementAction, err := parser.ParseSql(sql)

	if err != nil {
		t.Fatalf("Failed to parse SQL: %v", err)
	}
	if statementAction.Target != "name" || statementAction.Action != parser.SqlActionRenameColumn || statementAction.Relation != "rules" {
		t.Fatalf("Expected target 'name', action 'RENAME', relation 'rules', got target '%s', action '%s', relation '%s'", statementAction.Target, statementAction.Action, statementAction.Relation)
	}
}
