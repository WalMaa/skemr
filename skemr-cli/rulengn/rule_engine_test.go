package rulengn

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/walmaa/skemr-common/models"
)

func TestProcessStatementsReturnsResult(t *testing.T) {
	ruleEngine := NewRuleEngine()
	relName := "users"
	rules := []models.Rule{

		{
			ID:           uuid.New(),
			Name:         "Drop Age Column Rule",
			Type:         models.RuleTypeDeprecated,
			Scope:        models.RuleScopeColumn,
			RelationName: &relName,
			Target:       "age",
		},
	}
	statementDtos := []SqlStatementDto{
		{
			Statement: "ALTER TABLE users DROP COLUMN age;",
			File:      "test.sql",
			Line:      1,
		},
		{
			Statement: "ALTER TABLE users DROP COLUMN firstName;",
			File:      "test.sql",
			Line:      2,
		},
	}

	result, err := ruleEngine.ProcessStatements(t.Context(), statementDtos, rules)

	assert.NoError(t, err)

	assert.Equal(t, 1, len(result))
	assert.Equal(t, models.RuleTypeDeprecated, result[0].Type)
	assert.Equal(t, "Drop Age Column Rule", result[0].Rule.Name)
	assert.Equal(t, statementDtos[0].Statement, result[0].Statement)
	assert.Equal(t, "test.sql", result[0].File)
	assert.Equal(t, 1, result[0].Line)

}

func TestDeprecatedRuleTrigger(t *testing.T) {
	ruleEngine := NewRuleEngine()

	// Define a sample SQL statement
	statement := "ALTER TABLE users DROP COLUMN age;"
	statementDto := SqlStatementDto{
		Statement: statement,
		File:      "test.sql",
		Line:      1,
	}

	relName := "users"

	// Define a rules slice
	rules := []models.Rule{
		{
			ID:           uuid.New(),
			Name:         "Drop Age Column Rule",
			Type:         models.RuleTypeDeprecated,
			Scope:        models.RuleScopeColumn,
			RelationName: &relName,
			Target:       "age",
		},
	}

	// Call the CheckStatement method
	result, err := ruleEngine.CheckStatement(statementDto, rules)

	assert.NoError(t, err)

	// Assert the result
	assert.Equal(t, 1, len(result))
	assert.Equal(t, models.RuleTypeDeprecated, result[0].Type)
	assert.Equal(t, "Drop Age Column Rule", result[0].Rule.Name)
	assert.Equal(t, statement, result[0].Statement)
	assert.Equal(t, "test.sql", result[0].File)
	assert.Equal(t, 1, result[0].Line)
}

func TestWarnRuleTrigger(t *testing.T) {
	ruleEngine := NewRuleEngine()

	// Define a sample SQL statement
	statement := "ALTER TABLE users DROP COLUMN age;"
	statementDto := SqlStatementDto{
		Statement: statement,
		File:      "test.sql",
		Line:      1,
	}

	relName := "users"

	// Define a rules slice
	rules := []models.Rule{
		{
			ID:           uuid.New(),
			Name:         "Drop Age Column Rule",
			Type:         models.RuleTypeWarn,
			Scope:        models.RuleScopeColumn,
			RelationName: &relName,
			Target:       "age",
		},
	}

	// Call the CheckStatement method
	result, err := ruleEngine.CheckStatement(statementDto, rules)

	assert.NoError(t, err)

	// Assert the result
	assert.Equal(t, 1, len(result))
	assert.Equal(t, models.RuleTypeWarn, result[0].Type)
	assert.Equal(t, "Drop Age Column Rule", result[0].Rule.Name)
	assert.Equal(t, statement, result[0].Statement)
	assert.Equal(t, "test.sql", result[0].File)
	assert.Equal(t, 1, result[0].Line)
}

func TestLockedRuleViolation(t *testing.T) {
	ruleEngine := NewRuleEngine()

	// Define a sample SQL statement
	statement := "ALTER TABLE users DROP COLUMN age;"
	statementDto := SqlStatementDto{
		Statement: statement,
		File:      "test.sql",
		Line:      1,
	}

	relName := "users"

	// Define a rules slice
	rules := []models.Rule{
		{
			ID:           uuid.New(),
			Name:         "Drop Age Column Rule",
			Type:         models.RuleTypeLocked,
			Scope:        models.RuleScopeColumn,
			RelationName: &relName,
			Target:       "age",
		},
	}

	// Call the CheckStatement method
	result, err := ruleEngine.CheckStatement(statementDto, rules)

	assert.NoError(t, err)

	// Assert the result
	assert.Equal(t, 1, len(result))
	assert.Equal(t, models.RuleTypeLocked, result[0].Type)
	assert.Equal(t, "Drop Age Column Rule", result[0].Rule.Name)
	assert.Equal(t, statement, result[0].Statement)
	assert.Equal(t, "test.sql", result[0].File)
	assert.Equal(t, 1, result[0].Line)
}

func TestAdvisoryRuleTrigger(t *testing.T) {
	ruleEngine := NewRuleEngine()

	// Define a sample SQL statement
	statement := "ALTER TABLE users DROP COLUMN age;"
	statementDto := SqlStatementDto{
		Statement: statement,
		File:      "test.sql",
		Line:      1,
	}

	relName := "users"

	// Define a rules slice
	rules := []models.Rule{
		{
			ID:           uuid.New(),
			Name:         "Drop Age Column Rule",
			Type:         models.RuleTypeAdvisory,
			Scope:        models.RuleScopeColumn,
			RelationName: &relName,
			Target:       "age",
		},
	}

	// Call the CheckStatement method
	result, err := ruleEngine.CheckStatement(statementDto, rules)

	assert.NoError(t, err)

	// Assert the result
	assert.Equal(t, 1, len(result))
}
