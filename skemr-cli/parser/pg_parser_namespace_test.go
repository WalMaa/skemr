package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSqlDropSchema(t *testing.T) {
	sql := "DROP SCHEMA public"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "public", statementAction[0].Target, "Expected target 'public'")
	assert.Equal(t, SqlActionDropNamespace, statementAction[0].Action, "Expected action 'DROP SCHEMA'")
	assert.Equal(t, "", statementAction[0].Relation, "Expected empty relation for DROP SCHEMA")
}

func TestParseSqlCreateSchema(t *testing.T) {
	sql := "CREATE SCHEMA public"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "public", statementAction[0].Target, "Expected target 'public'")
	assert.Equal(t, SqlActionCreateNamespace, statementAction[0].Action, "Expected action 'CREATE SCHEMA'")
	assert.Equal(t, "", statementAction[0].Relation, "Expected empty relation for CREATE SCHEMA")
}

func TestParseSqlRenameSchema(t *testing.T) {
	sql := "ALTER SCHEMA public RENAME TO new_public"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "public", statementAction[0].Target, "Expected target 'public'")
	assert.Equal(t, SqlActionRenameNamespace, statementAction[0].Action, "Expected action 'RENAME SCHEMA'")
	assert.Equal(t, "", statementAction[0].Relation, "Expected empty relation for RENAME SCHEMA")
}
