package models

import (
	"time"

	"github.com/google/uuid"
)

type RuleAttributes struct {
	RemovalDate *string `json:"removalDate" validate:"omitempty,date_format=2006-01-02T15:04:05Z07:00"`
}

type Rule struct {
	ID             uuid.UUID      `json:"id"`
	Name           string         `json:"name"`
	Attributes     RuleAttributes `json:"attributes"`
	RuleType       RuleType       `json:"ruleType"`
	DataBaseEntity DatabaseEntity `json:"databaseEntity"`
	CreatedAt      time.Time      `json:"createdAt"`
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
