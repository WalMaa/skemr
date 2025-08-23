package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseSqlDropColumn(t *testing.T) {
	sql := "ALTER TABLE rules DROP COLUMN name"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "name", statementAction.Target, "Expected target 'name'")
	assert.Equal(t, SqlActionDropColumn, statementAction.Action, "Expected action 'DROP COLUMN'")
	assert.Equal(t, "rules", statementAction.Relation, "Expected relation 'rules'")

}

func TestParseSqlDropTable(t *testing.T) {
	sql := "DROP TABLE rules"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "", statementAction.Target, "Expected empty target for DROP TABLE")
	assert.Equal(t, SqlActionDropTable, statementAction.Action, "Expected action 'DROP TABLE'")
	assert.Equal(t, "rules", statementAction.Relation, "Expected relation 'rules'")
}

func TestParseSqlRenameColumn(t *testing.T) {
	sql := "ALTER TABLE rules RENAME COLUMN name TO new_name"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "name", statementAction.Target, "Expected target 'new_name'")
	if statementAction.Target != "name" || statementAction.Action != SqlActionRenameColumn || statementAction.Relation != "rules" {
		t.Fatalf("Expected target 'name', action 'RENAME', relation 'rules', got target '%s', action '%s', relation '%s'", statementAction.Target, statementAction.Action, statementAction.Relation)
	}
}

func TestParseSqlModifyDataType(t *testing.T) {
	sql := "ALTER TABLE rules ALTER COLUMN name TYPE VARCHAR(255)"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	if statementAction.Target != "name" || statementAction.Action != SqlActionModifyDataType || statementAction.Relation != "rules" {
		t.Fatalf("Expected target 'name', action 'MODIFY DATA TYPE', relation 'rules', got target '%s', action '%s', relation '%s'", statementAction.Target, statementAction.Action, statementAction.Relation)
	}
}

func TestParseSqlRenameTable(t *testing.T) {
	sql := "ALTER TABLE rules RENAME TO new_rules"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "", statementAction.Target, "Expected empty target for RENAME TABLE")
	assert.Equal(t, SqlActionRenameTable, statementAction.Action, "Expected action 'RENAME TABLE'")
	assert.Equal(t, "rules", statementAction.Relation, "Expected relation 'rules'")
}

func TestParseSqlDropDataBase(t *testing.T) {
	sql := "DROP DATABASE postgres"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	if statementAction.Target != "postgres" || statementAction.Action != SqlActionDropDatabase || statementAction.Relation != "" {
		t.Fatalf("Expected target 'postgres', action 'DROP DATABASE', relation '', got target '%s', action '%s', relation '%s'", statementAction.Target, statementAction.Action, statementAction.Relation)
	}
}

func TestParseSqlUndefined(t *testing.T) {
	sql := "This is not a valid SQL statement"
	statementAction, err := ParseSql(sql)

	if err == nil {
		t.Fatalf("Expected an error for invalid SQL, got none")
	}
	if statementAction.Target != "" || statementAction.Action != SqlActionUndefined || statementAction.Relation != "" {
		t.Fatalf("Expected target '', action 'UNDEFINED', relation '', got target '%s', action '%s', relation '%s'", statementAction.Target, statementAction.Action, statementAction.Relation)
	}
}
