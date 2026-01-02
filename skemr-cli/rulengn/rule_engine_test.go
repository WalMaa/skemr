package rulengn

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/walmaa/skemr-common/models"
)

func TestLockedRuleViolationExitCode1(t *testing.T) {
	ruleEngine := NewRuleEngine()

	// Define a sample SQL statement
	statement := "ALTER TABLE users DROP COLUMN age;"

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
	result := ruleEngine.CheckStatement(statement, rules)

	// Assert the result
	assert.True(t, result)
}

func TestCheckStatement(t *testing.T) {
	ruleEngine := NewRuleEngine()

	// Define a sample SQL statement
	statement := "ALTER TABLE users DROP COLUMN age;"

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
	result := ruleEngine.CheckStatement(statement, rules)

	// Assert the result
	assert.True(t, result)
}

func TestProcessStatement(t *testing.T) {
	ruleEngine := NewRuleEngine()

	// Define a sample SQL statement
	statements := []string{"ALTER TABLE users DROP COLUMN age;",
		"CREATE TABLE orders (id SERIAL PRIMARY KEY, user_id INT, amount DECIMAL);",
		"INSERT INTO users (name, email) VALUES ('John Doe', 'johg');"}

	// Call the processStatements method
	results, err := ruleEngine.ProcessStatements(context.Background(), statements, nil)

	assert.NoError(t, err)

	// Consume the results to trigger execution and logs
	var resultCount int
	for res := range results {
		assert.True(t, res)
		resultCount++
	}

	// Assert we received the expected number of results
	assert.Equal(t, len(statements), resultCount)

}
