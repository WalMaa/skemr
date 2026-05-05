package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSqlDropDataBase(t *testing.T) {
	sql := "DROP DATABASE postgres"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "postgres", statementAction[0].Target, "Expected target 'postgres'")
	assert.Equal(t, SqlActionDropDatabase, statementAction[0].Action, "Expected action 'SqlActionDropDatabase'")
	assert.Equal(t, "", statementAction[0].Relation)
}

func TestParseSqlCreateDataBase(t *testing.T) {
	sql := "CREATE DATABASE skemr_db"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "skemr_db", statementAction[0].Target, "Expected target 'skemr_db'")
	assert.Equal(t, SqlActionCreateDatabase, statementAction[0].Action, "Expected action 'CREATE DATABASE'")
	assert.Equal(t, "", statementAction[0].Relation, "Expected empty relation for CREATE DATABASE")
}

func TestParseSqlRenameDataBase(t *testing.T) {
	sql := "ALTER DATABASE skemr_db RENAME TO skemr_database"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "skemr_db", statementAction[0].Target, "Expected target 'skemr_db'")
	assert.Equal(t, SqlActionRenameDatabase, statementAction[0].Action, "Expected action 'RENAME DATABASE'")
	assert.Equal(t, "", statementAction[0].Relation, "Expected empty relation for RENAME DATABASE")
}
