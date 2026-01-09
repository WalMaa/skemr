package mapper

import (
	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/dto"
	"github.com/walmaa/skemr-common/models"
)

func ToDomainRule(e sqlc.Rule) models.Rule {
	return models.Rule{
		ID:               e.ID,
		Name:             e.Name,
		RuleType:         models.RuleType(e.Type),
		DataBaseEntityId: e.DatabaseEntityID,
		ProjectId:        e.ProjectID,
	}
}

func ToDomainRules(r []sqlc.Rule) []models.Rule {
	rules := make([]models.Rule, len(r))
	for i, rule := range r {
		rules[i] = ToDomainRule(rule)
	}
	return rules
}

func ToSqlcCreateRule(projectId uuid.UUID, dto dto.RuleCreationDto) sqlc.CreateRuleParams {
	return sqlc.CreateRuleParams{
		Name:      dto.Name,
		Type:      sqlc.RuleType(dto.RuleType),
		ProjectID: projectId,
	}
}
