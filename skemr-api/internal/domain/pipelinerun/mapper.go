package pipelinerun

import (
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/mapper"
	"github.com/walmaa/skemr-common/models"
)

func ToDomainPipelineRun(r sqlc.PipelineRun) models.PipelineRun {
	return models.PipelineRun{
		ID:          r.ID,
		Status:      models.MigrationStatus(r.Status),
		Environment: r.Environment.String,
		StartedAt:   mapper.Time(&r.StartedAt),
		CompletedAt: mapper.Time(&r.CompletedAt),
		CreatedAt:   mapper.Time(&r.CreatedAt),
	}
}

func ToDomainPipelineRuns(r []sqlc.PipelineRun) []models.PipelineRun {
	pipelineRuns := make([]models.PipelineRun, len(r))
	for i, run := range r {
		pipelineRuns[i] = ToDomainPipelineRun(run)
	}
	return pipelineRuns
}
