package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/mapper"
	"github.com/walmaa/skemr-common/models"
)

type RuleService struct {
	db sqlc.Querier
}

func NewRuleService(q sqlc.Querier) *RuleService {
	return &RuleService{db: q}
}

func (r *RuleService) CreateRule(c context.Context, projectID uuid.UUID, dto models.RuleCreationDto) (models.Rule, error) {
	slog.Info("Creating rule")

	_, err := CheckProjectExists(c, r.db, projectID)

	if err != nil {
		return models.Rule{}, err
	}

	rule, err := r.db.CreateRule(c, mapper.ToSqlcCreate(projectID, dto))
	if err != nil {
		slog.Error("Unable to create a Rule")
		return models.Rule{}, err
	}

	return mapper.ToDomainEntity(rule), nil
}
