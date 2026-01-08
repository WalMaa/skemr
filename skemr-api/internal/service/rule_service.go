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

func (r *RuleService) CreateRule(c context.Context, projectID uuid.UUID, databaseId uuid.UUID, dto models.RuleCreationDto) (models.Rule, error) {
	slog.Info("Creating rule")

	project, err := CheckProjectExists(c, r.db, projectID)

	if err != nil {
		return models.Rule{}, err
	}

	_, err = r.db.GetDatabaseByIdAndProject(c, sqlc.GetDatabaseByIdAndProjectParams{
		ID:        databaseId,
		ProjectID: project.ID,
	})

	if err != nil {
		slog.Error("Error fetching database", err)
		return models.Rule{}, err
	}

	rule, err := r.db.CreateRule(c, mapper.ToSqlcCreateRule(projectID, dto))
	if err != nil {
		slog.Error("Unable to create a Rule")
		return models.Rule{}, err
	}

	return mapper.ToDomainRule(rule), nil
}

func (r *RuleService) ListRulesByDatabase(c context.Context, projectID uuid.UUID, databaseID uuid.UUID) ([]models.Rule, error) {
	slog.Info("Listing rules for project %q and database %q", projectID, databaseID)

	project, err := CheckProjectExists(c, r.db, projectID)

	if err != nil {
		return []models.Rule{}, err
	}

	database, err := r.db.GetDatabaseByIdAndProject(c, sqlc.GetDatabaseByIdAndProjectParams{
		ID:        databaseID,
		ProjectID: project.ID,
	})

	if err != nil {
		slog.Error("Error fetching database", err)
		return []models.Rule{}, err
	}

	rules, err := r.db.ListRulesByDatabaseId(c, database.ID)
	return mapper.ToDomainRules(rules), nil

}
