package mapper

import (
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/dto"
	"github.com/walmaa/skemr-common/models"
)

func ToDomainRule(e sqlc.Rule) models.Rule {
	return models.Rule{
		ID:         e.ID,
		Name:       e.Name,
		RuleType:   models.RuleType(e.Type),
		Attributes: ToRuleAttributes(e.Attributes),
		CreatedAt:  Time(&e.CreatedAt),
	}
}

func ToRuleAttributes(attributes []byte) models.RuleAttributes {
	var ruleAttributes models.RuleAttributes
	err := json.Unmarshal(attributes, &ruleAttributes)
	if err != nil {
		slog.Error("Unable to unmarshal rule attributes", "error", err)
		panic(err)
	}
	return ruleAttributes
}

func ToDomainRuleWithEntity(e sqlc.GetRuleWithEntityRow) models.Rule {
	return models.Rule{
		ID:             e.Rule.ID,
		Name:           e.Rule.Name,
		Attributes:     ToRuleAttributes(e.Rule.Attributes),
		RuleType:       models.RuleType(e.Rule.Type),
		DataBaseEntity: ToDomainDatabaseEntity(e.DatabaseEntity),
		CreatedAt:      Time(&e.Rule.CreatedAt),
	}
}

func ToDomainRules(r []sqlc.Rule) []models.Rule {
	rules := make([]models.Rule, len(r))
	for i, rule := range r {
		rules[i] = ToDomainRule(rule)
	}
	return rules
}

func ToDomainRulesWithEntity(r []sqlc.GetRulesWithEntitiesRow) []models.Rule {
	rules := make([]models.Rule, len(r))
	for i, rule := range r {
		rules[i] = ToDomainRuleWithEntity(sqlc.GetRuleWithEntityRow(rule))
	}
	return rules
}

func ToSqlcCreateRule(databaseId uuid.UUID, dto dto.RuleCreationDto) sqlc.CreateRuleParams {
	return sqlc.CreateRuleParams{
		Name:             dto.Name,
		Type:             sqlc.RuleType(dto.RuleType),
		DatabaseID:       databaseId,
		Attributes:       ToBytes(dto.Attributes),
		DatabaseEntityID: dto.DataBaseEntityId,
	}
}
