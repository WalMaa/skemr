package pipelinerun

import "github.com/walmaa/skemr-common/models"

type PipelineRunCreationDto struct {
	Status      models.MigrationStatus `json:"status" validate:"required,oneof=completed failed"`
	Environment string                 `json:"environment" validate:"required"`
	CompletedAt string                 `json:"completedAt" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
}
