package models

import (
	"github.com/google/uuid"
)

type Rule struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Type         RuleType  `json:"type"`
	Scope        RuleScope `json:"scope"`
	RelationName *string   `json:"relation_name"`
	Target       string    `json:"target"`
	DatabaseID   uuid.UUID `json:"database_id"`
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
