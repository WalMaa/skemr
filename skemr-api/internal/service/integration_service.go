package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/walmaa/skemr-common/models"
)

type IntegrationService struct {
	RuleService *RuleService
}

func NewIntegrationService(ruleService *RuleService) *IntegrationService {
	return &IntegrationService{RuleService: ruleService}
}

func (s *IntegrationService) ListRulesByDatabase(c context.Context, projectID uuid.UUID, databaseID uuid.UUID) ([]models.Rule, error) {
	return s.RuleService.ListRulesByDatabase(c, projectID, databaseID)
}
