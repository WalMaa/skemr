package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, SqlActionRenameColumn, statementAction.Action, "Expected action 'RENAME COLUMN'")
	assert.Equal(t, "rules", statementAction.Relation, "Expected relation 'rules'")
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

	assert.Equal(t, "rules", statementAction.Target, "Expected target 'rules'")
	assert.Equal(t, SqlActionRenameTable, statementAction.Action, "Expected action 'RENAME TABLE'")
	assert.Equal(t, "", statementAction.Relation, "Expected empty relation for RENAME TABLE")
}

func TestParseSqlDropDataBase(t *testing.T) {
	sql := "DROP DATABASE postgres"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	if statementAction.Target != "postgres" || statementAction.Action != SqlActionDropDatabase || statementAction.Relation != "" {
		t.Fatalf("Expected target 'postgres', action 'DROP DATABASE', relation '', got target '%s', action '%s', relation '%s'", statementAction.Target, statementAction.Action, statementAction.Relation)
	}
}

func TestParseSqlCreateDataBase(t *testing.T) {
	sql := "CREATE DATABASE skemr_db"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "skemr_db", statementAction.Target, "Expected target 'skemr_db'")
	assert.Equal(t, SqlActionCreateDatabase, statementAction.Action, "Expected action 'CREATE DATABASE'")
	assert.Equal(t, "", statementAction.Relation, "Expected empty relation for CREATE DATABASE")
}

func TestParseSqlRenameDataBase(t *testing.T) {
	sql := "ALTER DATABASE skemr_db RENAME TO skemr_database"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "skemr_db", statementAction.Target, "Expected target 'skemr_db'")
	assert.Equal(t, SqlActionRenameDatabase, statementAction.Action, "Expected action 'RENAME DATABASE'")
	assert.Equal(t, "", statementAction.Relation, "Expected empty relation for RENAME DATABASE")
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

func TestParseSqlAddColumn(t *testing.T) {
	sql := "ALTER TABLE rules ADD COLUMN description TEXT"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "description", statementAction.Target, "Expected target 'description'")
	assert.Equal(t, SqlActionAddColumn, statementAction.Action, "Expected action 'ADD COLUMN'")
	assert.Equal(t, "rules", statementAction.Relation, "Expected relation 'rules'")
}

func TestParseSqlInsertRow(t *testing.T) {
	sql := "INSERT INTO rules (name, scope) VALUES ('rule1', 'table')"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "", statementAction.Target, "Expected empty target for INSERT ROW")
	assert.Equal(t, SqlActionInsertRow, statementAction.Action, "Expected action 'INSERT ROW'")
	assert.Equal(t, "rules", statementAction.Relation, "Expected relation 'rules'")
}
