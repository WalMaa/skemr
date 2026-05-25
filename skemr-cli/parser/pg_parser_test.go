package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCreateIndex(t *testing.T) {
	sql := "CREATE INDEX idx_name ON rules (name)"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "", statementAction[0].Target, "Expected target ''")
	assert.Equal(t, SqlActionUndefined, statementAction[0].Action, "Expected action 'UNDEFINED' for unsupported statement")
	assert.Equal(t, "", statementAction[0].Relation, "Expected relation ''")
}

func TestParseSqlUndefined(t *testing.T) {
	sql := "This is not a valid SQL statement"
	statementAction, err := ParseSql(sql)

	assert.NotNil(t, err, "Expected an error for invalid SQL, got none")
	assert.Nil(t, statementAction, "Expected statementAction to be nil for invalid SQL")
}

func TestParseSqlInsertRow(t *testing.T) {
	sql := "INSERT INTO rules (name, scope) VALUES ('rule1', 'table')"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)

	assert.Equal(t, "", statementAction[0].Target, "Expected empty target for INSERT ROW")
	assert.Equal(t, SqlActionInsertRow, statementAction[0].Action, "Expected action 'INSERT ROW'")
	assert.Equal(t, "rules", statementAction[0].Relation, "Expected relation 'rules'")
}

func TestParseMultipleStatements(t *testing.T) {
	sql := "ALTER TABLE rules DROP COLUMN name;      \n" +
		"ALTER TABLE users DROP COLUMN createdAt;"

	statementActions, err := ParseSql(sql)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(statementActions))
}

func TestParseMultipleStatementsWithComments(t *testing.T) {
	sql := "ALTER TABLE rules DROP COLUMN name;      \n" +
		"-- This is a comment \n" +
		"ALTER TABLE users DROP COLUMN createdAt;" +
		"ALTER TABLE \n" +
		"orders\n" +
		"ADD COLUMN order_date DATE;"

	statementActions, err := ParseSql(sql)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(statementActions))
}

func TestParserOriginalString(t *testing.T) {
	sql := "ALTER TABLE rules DROP COLUMN name;"
	statementAction, err := ParseSql(sql)

	expected := strings.Replace(sql, ";", "", -1)
	assert.Nil(t, err)
	assert.Equal(t, expected, statementAction[0].Original, "Expected original string to match input SQL")
}

func TestParserOriginalStringNoSemicolon(t *testing.T) {
	sql := "ALTER TABLE rules DROP COLUMN name"
	statementAction, err := ParseSql(sql)

	assert.Nil(t, err)
	assert.Equal(t, sql, statementAction[0].Original, "Expected original string to match input SQL")
}

func TestParserOriginalStringNewlines(t *testing.T) {
	sql := "ALTER TABLE \n " +
		"rules\r\n" +
		"DROP COLUMN name;"
	statementAction, err := ParseSql(sql)

	expected := "ALTER TABLE rules DROP COLUMN name"
	assert.Nil(t, err)
	assert.Equal(t, expected, statementAction[0].Original, "Expected original string to have newlines replaced with spaces and semicolon removed")
}

func TestParserOriginalStringTrimSpaces(t *testing.T) {
	sql := "   ALTER TABLE rules DROP COLUMN name;   "
	statementAction, err := ParseSql(sql)

	expected := "ALTER TABLE rules DROP COLUMN name"
	assert.Nil(t, err)
	assert.Equal(t, expected, statementAction[0].Original, "Expected original string to be trimmed of leading and trailing spaces and semicolon removed")
}
