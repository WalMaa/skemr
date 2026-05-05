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
