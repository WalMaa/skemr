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
		Target: "age",
	}
	sql := "ALTER TABLE users DROP COLUMN age"
	_, _ = parser.ParseRule(rule, sql)

}
