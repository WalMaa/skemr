package rules

import (
	"github.com/google/uuid"
	"github.com/walmaa/skemr-common/models"
)

type RuleCreationDto struct {
	Name             string
	RuleType         models.RuleType       `json:"ruleType" validate:"required,oneof=locked deprecated advisory warning"`
	Attributes       models.RuleAttributes `json:"attributes" validate:"omitempty,json"`
	DataBaseEntityId uuid.UUID             `json:"databaseEntityId" validate:"required,uuid4"`
}
