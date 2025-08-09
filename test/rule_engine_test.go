package test

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/walmaa/skemr/db/sqlc"
	"github.com/walmaa/skemr/rulengn"
	"github.com/walmaa/skemr/test/mocks"
	"testing"
)

func TestCheckStatement(t *testing.T) {
	// Initialize the rule engine with a mock database
	mockDB := &mocks.Querier{}
	ruleEngine := rulengn.NewRuleEngine(mockDB)

	// Define a sample project
	database := &sqlc.Database{ID: uuid.New(), Name: "Test Database"}

	// Define a sample SQL statement
	statement := "ALTER TABLE users DROP COLUMN age;"
	relation := "users"
	// Mock the database response
	mockDB.On("ListRulesByCriteria", mock.Anything, sqlc.ListRulesByCriteriaParams{
		DatabaseID:   database.ID,
		Scope:        sqlc.RuleScopeTable,
		RelationName: &relation,
		Target:       "age",
	}).Return([]sqlc.Rule{}, nil)

	// Call the CheckStatement method
	result := ruleEngine.CheckStatement(context.Background(), statement, database)

	// Assert the result
	assert.True(t, result)
}

func TestProcessStatement(t *testing.T) {
	// Initialize the rule engine with a mock database
	mockDB := &mocks.Querier{}
	ruleEngine := rulengn.NewRuleEngine(mockDB)

	// Define a sample project
	database := &sqlc.Database{ID: uuid.New(), Name: "Test Database"}

	// Define a sample SQL statement
	statements := []string{"ALTER TABLE users DROP COLUMN age;",
		"CREATE TABLE orders (id SERIAL PRIMARY KEY, user_id INT, amount DECIMAL);",
		"INSERT INTO users (name, email) VALUES ('John Doe', 'johg');"}

	// Mock the database response
	mockDB.On("ListRulesByCriteria", mock.Anything, mock.Anything).Return([]sqlc.Rule{}, nil)

	// Call the processStatements method
	results, err := ruleEngine.ProcessStatements(context.Background(), statements, database)

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
