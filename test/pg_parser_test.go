package test

import (
	"github.com/walmaa/skemr/parser"
	"testing"
)

//func TestSQLParser(t *testing.T) {
//	rule := &skemr.Rule{
//		ID:     1,
//		Name:   "Test Rule",
//		Type:   skemr.RuleTypeLock,
//		Scope:  skemr.RuleScopeColumn,
//		Target: "name",
//	}
//	sql := "ALTER TABLE rules DROP COLUMN name"
//	_, _ = parser.ParseRule(rule, sql)
//
//}

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

func TestParseSqlModifyDataType(t *testing.T) {
	sql := "ALTER TABLE rules ALTER COLUMN name TYPE VARCHAR(255)"
	statementAction, err := parser.ParseSql(sql)

	if err != nil {
		t.Fatalf("Failed to parse SQL: %v", err)
	}
	if statementAction.Target != "name" || statementAction.Action != parser.SqlActionModifyDataType || statementAction.Relation != "rules" {
		t.Fatalf("Expected target 'name', action 'MODIFY DATA TYPE', relation 'rules', got target '%s', action '%s', relation '%s'", statementAction.Target, statementAction.Action, statementAction.Relation)
	}
}

func TestParseSqlDropDataBase(t *testing.T) {
	sql := "DROP DATABASE postgres"
	statementAction, err := parser.ParseSql(sql)

	if err != nil {
		t.Fatalf("Failed to parse SQL: %v", err)
	}
	if statementAction.Target != "postgres" || statementAction.Action != parser.SqlActionDropDatabase || statementAction.Relation != "" {
		t.Fatalf("Expected target 'postgres', action 'DROP DATABASE', relation '', got target '%s', action '%s', relation '%s'", statementAction.Target, statementAction.Action, statementAction.Relation)
	}
}

func TestParseSqlUndefined(t *testing.T) {
	sql := "This is not a valid SQL statement"
	statementAction, err := parser.ParseSql(sql)

	if err == nil {
		t.Fatalf("Expected an error for invalid SQL, got none")
	}
	if statementAction.Target != "" || statementAction.Action != parser.SqlActionUndefined || statementAction.Relation != "" {
		t.Fatalf("Expected target '', action 'UNDEFINED', relation '', got target '%s', action '%s', relation '%s'", statementAction.Target, statementAction.Action, statementAction.Relation)
	}
}
