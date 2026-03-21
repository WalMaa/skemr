package rulengn

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/walmaa/skemr-common/models"
)

var entities = []models.DatabaseEntity{
	{
		ID:   uuid.New(),
		Name: "users",
		Type: models.DatabaseEntityTypeTable,
	},
	{
		ID:       uuid.New(),
		Name:     "age",
		Type:     models.DatabaseEntityTypeColumn,
	},
	{
		ID:   uuid.New(),
		Name: "customers",
		Type: models.DatabaseEntityTypeTable,
	},
	{
		ID:       uuid.New(),
		Name:     "age",
		Type:     models.DatabaseEntityTypeColumn,
	},
}

func TestProcessStatementsReturnsResult(t *testing.T) {
	ruleEngine := NewRuleEngine()

	// Use existing test migration file
	migrationFile := filepath.Join("..", "test", "sql", "migration.sql")

	rules := []models.Rule{
		{
			ID:       uuid.New(),
			Name:     "Drop Column Rule",
			RuleType: models.RuleTypeDeprecated,
			DataBaseEntity: models.DatabaseEntity{
				Name: "password_hash", // matches a DROP COLUMN in migration.sql
			},
		},
	}
	statementDtos := []MigrationFileDto{
		{File: migrationFile},
	}

	result, err := ruleEngine.ProcessMigrationFiles(t.Context(), statementDtos, rules, entities)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(result), 1)
	assert.Equal(t, models.RuleTypeDeprecated, result[0].Type)
	assert.Equal(t, "Drop Column Rule", result[0].Rule.Name)
	assert.Equal(t, migrationFile, result[0].File)
}

func TestDeprecatedRuleTrigger(t *testing.T) {
	ruleEngine := NewRuleEngine()

	// Create a temporary migration file with a statement that drops column age
	tmpFile, err := os.CreateTemp(t.TempDir(), "migration-*.sql")
	assert.NoError(t, err)
	defer func() { _ = tmpFile.Close() }()
	content := "CREATE TABLE users (age INT);\nALTER TABLE users DROP COLUMN age;"
	_, err = tmpFile.WriteString(content)
	assert.NoError(t, err)

	migrationFile := tmpFile.Name()

	rules := []models.Rule{
		{
			ID:       uuid.New(),
			Name:     "Drop Age Column Rule",
			RuleType: models.RuleTypeDeprecated,
			DataBaseEntity: models.DatabaseEntity{
				Name: "age",
			},
		},
	}

	result, err := ruleEngine.CheckStatement(MigrationFileDto{File: migrationFile}, rules, entities)
	assert.NoError(t, err)

	// Assert the result
	assert.Equal(t, 1, len(result))
	assert.Equal(t, models.RuleTypeDeprecated, result[0].Type)
	assert.Equal(t, "Drop Age Column Rule", result[0].Rule.Name)
	assert.Equal(t, migrationFile, result[0].File)
}

func TestWarnRuleTrigger(t *testing.T) {
	ruleEngine := NewRuleEngine()

	// Create a temporary migration file
	tmpFile, err := os.CreateTemp(t.TempDir(), "migration-*.sql")
	assert.NoError(t, err)
	defer func() { _ = tmpFile.Close() }()
	content := "CREATE TABLE users (age INT);\nALTER TABLE users DROP COLUMN age;"
	_, err = tmpFile.WriteString(content)
	assert.NoError(t, err)
	migrationFile := tmpFile.Name()

	rules := []models.Rule{
		{
			ID:       uuid.New(),
			Name:     "Drop Age Column Rule",
			RuleType: models.RuleTypeWarn,
			DataBaseEntity: models.DatabaseEntity{
				Name: "age",
			},
		},
	}

	result, err := ruleEngine.CheckStatement(MigrationFileDto{File: migrationFile}, rules, entities)
	assert.NoError(t, err)

	// Assert the result
	assert.Equal(t, 1, len(result))
	assert.Equal(t, models.RuleTypeWarn, result[0].Type)
	assert.Equal(t, "Drop Age Column Rule", result[0].Rule.Name)
	assert.Equal(t, migrationFile, result[0].File)
}

func TestLockedRuleViolation(t *testing.T) {
	ruleEngine := NewRuleEngine()

	// Create a temporary migration file
	tmpFile, err := os.CreateTemp(t.TempDir(), "migration-*.sql")
	assert.NoError(t, err)
	defer func() { _ = tmpFile.Close() }()
	content := "CREATE TABLE users (age INT);\nALTER TABLE users DROP COLUMN age;"
	_, err = tmpFile.WriteString(content)
	assert.NoError(t, err)
	migrationFile := tmpFile.Name()

	rules := []models.Rule{
		{
			ID:       uuid.New(),
			Name:     "Drop Age Column Rule",
			RuleType: models.RuleTypeLocked,
			DataBaseEntity: models.DatabaseEntity{
				Name: "age",
			},
		},
	}

	result, err := ruleEngine.CheckStatement(MigrationFileDto{File: migrationFile}, rules, entities)
	assert.NoError(t, err)

	// Assert the result
	assert.Equal(t, 1, len(result))
	assert.Equal(t, models.RuleTypeLocked, result[0].Type)
	assert.Equal(t, "Drop Age Column Rule", result[0].Rule.Name)
	assert.Equal(t, migrationFile, result[0].File)
}

func TestLockedTableRuleViolationOnColumnAdd(t *testing.T) {
	ruleEngine := NewRuleEngine()

	// Create a temporary migration file
	tmpFile, err := os.CreateTemp(t.TempDir(), "migration-*.sql")
	assert.NoError(t, err)
	defer func() { _ = tmpFile.Close() }()
	content := "CREATE TABLE users (age INT);\nALTER TABLE users ADD COLUMN name VARCHAR(255);"
	_, err = tmpFile.WriteString(content)
	assert.NoError(t, err)
	migrationFile := tmpFile.Name()

	rules := []models.Rule{
		{
			ID:       uuid.New(),
			Name:     "Locked Users Table Rule",
			RuleType: models.RuleTypeLocked,
			DataBaseEntity: models.DatabaseEntity{
				Name: "users",
			},
		},
	}

	result, err := ruleEngine.CheckStatement(MigrationFileDto{File: migrationFile}, rules, entities)
	assert.NoError(t, err)

	// Assert the result
	assert.Equal(t, 1, len(result))
	assert.Equal(t, models.RuleTypeLocked, result[0].Type)
	assert.Equal(t, "Locked Users Table Rule", result[0].Rule.Name)
	assert.Equal(t, migrationFile, result[0].File)
}

func TestAdvisoryRuleTrigger(t *testing.T) {
	ruleEngine := NewRuleEngine()

	// Create a temporary migration file
	tmpFile, err := os.CreateTemp(t.TempDir(), "migration-*.sql")
	assert.NoError(t, err)
	defer func() { _ = tmpFile.Close() }()
	content := "CREATE TABLE users (age INT);\nALTER TABLE users DROP COLUMN age;"
	_, err = tmpFile.WriteString(content)
	assert.NoError(t, err)
	migrationFile := tmpFile.Name()

	rules := []models.Rule{
		{
			ID:       uuid.New(),
			Name:     "Drop Age Column Rule",
			RuleType: models.RuleTypeAdvisory,
			DataBaseEntity: models.DatabaseEntity{
				Name: "age",
			},
		},
	}

	result, err := ruleEngine.CheckStatement(MigrationFileDto{File: migrationFile}, rules, entities)
	assert.NoError(t, err)

	// Assert the result
	assert.Equal(t, 1, len(result))
	assert.Equal(t, models.RuleTypeAdvisory, result[0].Type)
	assert.Equal(t, "Drop Age Column Rule", result[0].Rule.Name)
	assert.Equal(t, migrationFile, result[0].File)
}


// If tables A and B have identical column names, and there is a rule that locks on column name on table A,
// Then dropping the column on table B should not trigger the rule, but dropping the column on table A should trigger the rule.
func TestIdenticalColumnNameRule(t *testing.T) {
	tableAId := uuid.New()
	tableBId := uuid.New()
	tableAColumn := models.DatabaseEntity{
		ID:       uuid.New(),
		Name:     "age",
		Type:     models.DatabaseEntityTypeColumn,
		ParentId: &tableAId,
	}
	tableBColumn := models.DatabaseEntity{
		ID:       uuid.New(),
		Name:     "age",
		Type:     models.DatabaseEntityTypeColumn,
		ParentId: &tableBId,
	}

	entities := []models.DatabaseEntity{
		{
			ID:   tableAId,
			Name: "users",
			Type: models.DatabaseEntityTypeTable,
		},
		tableAColumn,
		{
			ID:   tableBId,
			Name: "customers",
			Type: models.DatabaseEntityTypeTable,
		},
		tableBColumn,
	}
	ruleEngine := NewRuleEngine()

	// Create a temporary migration file
	tmpFile, err := os.CreateTemp(t.TempDir(), "migration-*.sql")
	assert.NoError(t, err)
	defer func() { _ = tmpFile.Close() }()
	content := "ALTER TABLE customers DROP COLUMN age;\nALTER TABLE users DROP COLUMN age;"
	_, err = tmpFile.WriteString(content)
	assert.NoError(t, err)
	migrationFile := tmpFile.Name()

	rules := []models.Rule{
		{
			ID:       uuid.New(),
			Name:     "Locked Age Column on Users Table",
			RuleType: models.RuleTypeLocked,
			DataBaseEntity: tableAColumn,
		},
	}

	result, err := ruleEngine.CheckStatement(MigrationFileDto{File: migrationFile}, rules, entities)
	assert.NoError(t, err)

	// Assert the result
	assert.Equal(t, 1, len(result))
	assert.Equal(t, models.RuleTypeLocked, result[0].Type)
	assert.Equal(t, "Locked Age Column on Users Table", result[0].Rule.Name)
	assert.Equal(t, migrationFile, result[0].File)
}
