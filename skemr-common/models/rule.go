package models

import (
	"github.com/google/uuid"
)

type Rule struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Type             RuleType  `json:"type"`
	DataBaseEntityId uuid.UUID
	ProjectId        uuid.UUID
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

type RuleCreationDto struct {
	Name             string
	Type             RuleType
	DataBaseEntityId uuid.UUID
}
