package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/internal/domain/rules"
	"github.com/walmaa/skemr-api/internal/dto"
	"github.com/walmaa/skemr-common/models"
)

type IntegrationService struct {
	RuleService        *rules.RuleService
	PipelineRunService *PipelineRunService
}

func NewIntegrationService(ruleService *rules.RuleService, pipelineRunService *PipelineRunService) *IntegrationService {
	return &IntegrationService{RuleService: ruleService, PipelineRunService: pipelineRunService}
}

func (s *IntegrationService) ListRulesByDatabase(c context.Context, projectID uuid.UUID, databaseID uuid.UUID) ([]models.Rule, error) {
	return s.RuleService.ListRulesByDatabase(c, projectID, databaseID)
}

func (s *IntegrationService) CreatePipeLineRun(c context.Context, projectID uuid.UUID, databaseID uuid.UUID, dto dto.PipelineRunCreationDto) (models.PipelineRun, error) {
	return s.PipelineRunService.CreatePipelineRun(c, projectID, databaseID, dto)
}
