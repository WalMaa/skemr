package mapper

import (
	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-common/models"
)

func ToDomainEntity(e sqlc.Rule) models.Rule {
	return models.Rule{
		ID:               e.ID,
		Name:             e.Name,
		Type:             models.RuleType(e.Type),
		DataBaseEntityId: e.DatabaseEntityID,
		ProjectId:        e.ProjectID,
	}
}

func ToSqlcCreate(projectId uuid.UUID, dto models.RuleCreationDto) sqlc.CreateRuleParams {
	return sqlc.CreateRuleParams{
		Name:      dto.Name,
		Type:      sqlc.RuleType(dto.Type),
		ProjectID: projectId,
	}
}
