package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSqlDropColumn(t *testing.T) {
	sql := "ALTER TABLE rules DROP COLUMN name"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "name", statementAction[0].Target, "Expected target 'name'")
	assert.Equal(t, SqlActionDropColumn, statementAction[0].Action, "Expected action 'DROP COLUMN'")
	assert.Equal(t, "rules", statementAction[0].Relation, "Expected relation 'rules'")
}

func TestParseSqlRenameColumn(t *testing.T) {
	sql := "ALTER TABLE rules RENAME COLUMN name TO new_name"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "name", statementAction[0].Target, "Expected target 'name'")
	assert.Equal(t, SqlActionRenameColumn, statementAction[0].Action, "Expected action 'RENAME COLUMN'")
	assert.Equal(t, "rules", statementAction[0].Relation, "Expected relation 'rules'")
}

func TestParseDropQualifiedColumn(t *testing.T) {
	sql := "ALTER TABLE public.rules DROP COLUMN name"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "name", statementAction[0].Target, "Expected target 'name'")
	assert.Equal(t, SqlActionDropColumn, statementAction[0].Action, "Expected action 'DROP COLUMN'")
	assert.Equal(t, "rules", statementAction[0].Relation, "Expected relation 'rules'")
}

func TestParseSqlModifyColumnDataType(t *testing.T) {
	sql := "ALTER TABLE rules ALTER COLUMN name TYPE VARCHAR(255)"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "name", statementAction[0].Target, "Expected target 'name'")
	assert.Equal(t, SqlActionModifyDataType, statementAction[0].Action, "Expected action 'MODIFY DATA TYPE'")
	assert.Equal(t, "rules", statementAction[0].Relation, "Expected relation 'rules'")

}

func TestParseSqlAddColumn(t *testing.T) {
	sql := "ALTER TABLE rules ADD COLUMN description TEXT"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)
	assert.Equal(t, "description", statementAction[0].Target, "Expected target 'description'")
	assert.Equal(t, SqlActionAddColumn, statementAction[0].Action, "Expected action 'ADD COLUMN'")
	assert.Equal(t, "rules", statementAction[0].Relation, "Expected relation 'rules'")
}

func TestParseSqlAddColumnWithTimestamp(t *testing.T) {
	sql := "ALTER TABLE orders ADD COLUMN updated_at TIMESTAMP;"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)
	assert.Equal(t, "updated_at", statementAction[0].Target, "Expected target 'updated_at'")
	assert.Equal(t, SqlActionAddColumn, statementAction[0].Action, "Expected action 'ADD COLUMN'")
	assert.Equal(t, "orders", statementAction[0].Relation, "Expected relation 'orders'")
}
