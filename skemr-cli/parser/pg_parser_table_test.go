package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSqlDropTable(t *testing.T) {
	sql := "DROP TABLE rules"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "rules", statementAction[0].Target, "Expected empty target for DROP TABLE")
	assert.Equal(t, SqlActionDropTable, statementAction[0].Action, "Expected action 'DROP TABLE'")
	assert.Equal(t, "rules", statementAction[0].Relation, "Expected relation 'rules'")
}

func BenchmarkParseSqlDropTable(b *testing.B) {
	for b.Loop() {
		sql := "DROP TABLE rules"
		_, err := ParseSql(sql)
		assert.Nil(b, err)
	}
}

func TestParseSqlCreateTable(t *testing.T) {
	sql := "CREATE TABLE rules (id SERIAL PRIMARY KEY, name VARCHAR(255))"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "rules", statementAction[0].Target, "Expected target 'rules'")
	assert.Equal(t, SqlActionCreateTable, statementAction[0].Action, "Expected action 'CREATE TABLE'")
	assert.Equal(t, "rules", statementAction[0].Relation, "Expected relation 'rules'")
}

func TestParseSqlDropTableCascade(t *testing.T) {
	sql := "DROP TABLE rules CASCADE"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "rules", statementAction[0].Target, "Expected target 'rules'")
	assert.Equal(t, SqlActionDropTable, statementAction[0].Action, "Expected action 'DROP TABLE'")
	assert.Equal(t, "rules", statementAction[0].Relation, "Expected relation 'rules'")
}

func TestParseSqlRenameTable(t *testing.T) {
	sql := "ALTER TABLE rules RENAME TO new_rules"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "rules", statementAction[0].Target, "Expected target 'rules'")
	assert.Equal(t, SqlActionRenameTable, statementAction[0].Action, "Expected action 'RENAME TABLE'")
	assert.Equal(t, "", statementAction[0].Relation, "Expected empty relation for RENAME TABLE")
}

func TestParseDropQualifiedTable(t *testing.T) {
	sql := "DROP TABLE other.rules"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "rules", statementAction[0].Target, "Expected target 'rules'")
	assert.Equal(t, SqlActionDropTable, statementAction[0].Action, "Expected action 'DROP TABLE'")
	assert.Equal(t, "other", statementAction[0].Namespace, "Expected namespace 'other'")
}

func TestParseSqlRenameQualifiedTable(t *testing.T) {
	sql := "ALTER TABLE other.rules RENAME TO new_rules"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "rules", statementAction[0].Target, "Expected target 'rules'")
	assert.Equal(t, SqlActionRenameTable, statementAction[0].Action, "Expected action 'RENAME TABLE'")
	assert.Equal(t, "other", statementAction[0].Namespace, "Expected namespace 'other'")
}

func TestParseSqlCreateQualifiedTable(t *testing.T) {
	sql := "CREATE TABLE other.rules (id SERIAL PRIMARY KEY, name VARCHAR(255))"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "rules", statementAction[0].Target, "Expected target 'rules'")
	assert.Equal(t, SqlActionCreateTable, statementAction[0].Action, "Expected action 'CREATE TABLE'")
	assert.Equal(t, "other", statementAction[0].Namespace, "Expected namespace 'other'")
}
