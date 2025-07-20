package test

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"skemr/db/sqlc"
	"skemr/rulengn"
	"skemr/test/mocks"
	"testing"
)

func TestCheckStatement(t *testing.T) {
	// Initialize the rule engine with a mock database
	mockDB := &mocks.Querier{}
	ruleEngine := rulengn.NewRuleEngine(mockDB)

	// Define a sample project
	project := &sqlc.Project{ID: uuid.New(), Name: "Test Project"}

	// Define a sample SQL statement
	statement := "ALTER TABLE users DROP COLUMN age;"
	relation := "users"
	// Mock the database response
	mockDB.On("ListRulesByCriteria", mock.Anything, sqlc.ListRulesByCriteriaParams{
		ProjectID:    project.ID,
		Scope:        sqlc.RuleScopeTable,
		RelationName: &relation,
		Target:       "age",
	}).Return([]sqlc.Rule{}, nil)

	// Call the CheckStatement method
	result := ruleEngine.CheckStatement(context.Background(), statement, project)

	// Assert the result
	assert.True(t, result)
}
