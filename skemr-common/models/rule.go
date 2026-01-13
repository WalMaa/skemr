package models

import (
	"github.com/google/uuid"
)

type Rule struct {
	ID             uuid.UUID      `json:"id"`
	Name           string         `json:"name"`
	RuleType       RuleType       `json:"ruletype"`
	DataBaseEntity DatabaseEntity `json:"databaseEntity"`
}

type RuleType string

const (
	RuleTypeLocked     RuleType = "locked"
	RuleTypeWarn       RuleType = "warn"
	RuleTypeAdvisory   RuleType = "advisory"
	RuleTypeDeprecated RuleType = "deprecated"
)

type RuleScope string

const (
	RuleScopeDatabase RuleScope = "database"
	RuleScopeSchema   RuleScope = "schema"
	RuleScopeTable    RuleScope = "table"
	RuleScopeColumn   RuleScope = "column"
)
