package dto

import (
	"github.com/walmaa/skemr-common/models"
)

type SecretCreationDto struct {
	Name      string `json:"name" validate:"required,min=2,max=100"`
	ExpiresAt string `json:"expiresAt" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

type PipelineRunCreationDto struct {
	Status      models.MigrationStatus `json:"status" validate:"required,oneof=completed failed"`
	Environment string                 `json:"environment" validate:"required"`
	CompletedAt string                 `json:"completedAt" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
}
