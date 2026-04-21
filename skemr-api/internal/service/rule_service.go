package service

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/dto"
	"github.com/walmaa/skemr-api/internal/errormsg"
	"github.com/walmaa/skemr-api/internal/mapper"
	"github.com/walmaa/skemr-common/models"
)

type RuleService struct {
	db sqlc.Querier
}

func NewRuleService(q sqlc.Querier) *RuleService {
	return &RuleService{db: q}
}

func (r *RuleService) GetRule(c context.Context, projectID uuid.UUID, databaseID uuid.UUID, ruleID uuid.UUID) (models.Rule, error) {
	slog.Info("Fetching rule", "ruleID", ruleID)

	project, err := CheckProjectExists(c, r.db, projectID)

	if err != nil {
		return models.Rule{}, err
	}

	database, err := CheckDatabaseExists(c, r.db, project.ID, databaseID)

	if err != nil {
		slog.Error("Error fetching database", err)
		return models.Rule{}, err
	}

	rule, err := r.db.GetRuleWithEntity(c, sqlc.GetRuleWithEntityParams{
		DatabaseID: database.ID,
		RuleID:     ruleID,
	})

	if err != nil {
		slog.Error("Unable to fetch rule", "error", err)
		return models.Rule{}, err
	}

	return mapper.ToDomainRuleWithEntity(rule), nil
}

func (r *RuleService) CreateRule(c context.Context, projectID uuid.UUID, databaseId uuid.UUID, dto dto.RuleCreationDto) (models.Rule, error) {
	slog.Info("Creating rule")

	project, err := CheckProjectExists(c, r.db, projectID)

	if err != nil {
		return models.Rule{}, err
	}

	_, err = CheckDatabaseExists(c, r.db, project.ID, databaseId)

	if err != nil {
		slog.Error("Error fetching database", err)
		return models.Rule{}, err
	}

	// Check if a rule with the same name already exists
	exists, err := r.db.GetRuleByDatabaseAndName(c, sqlc.GetRuleByDatabaseAndNameParams{
		DatabaseID: databaseId,
		Name:       dto.Name,
	})

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("Error checking for existing rule", "name", dto.Name, "err", err)
		return models.Rule{}, err
	}

	if exists.ID != uuid.Nil {
		slog.Warn("Rule with the same name already exists", "name", dto.Name)
		return models.Rule{}, &models.ErrorResponse{
			Message: errormsg.ErrRuleWithSameName,
			Status:  http.StatusConflict,
		}
	}

	rule, err := r.db.CreateRule(c, mapper.ToSqlcCreateRule(databaseId, dto))
	if err != nil {
		slog.Error("Unable to create a Rule", err)
		return models.Rule{}, err
	}

	return mapper.ToDomainRule(rule), nil
}

func (r *RuleService) ListRulesByDatabase(c context.Context, projectID uuid.UUID, databaseID uuid.UUID) ([]models.Rule, error) {
	slog.Info("Listing rules", "projectID", projectID, "databaseID", databaseID)

	project, err := CheckProjectExists(c, r.db, projectID)

	if err != nil {
		return []models.Rule{}, err
	}

	database, err := CheckDatabaseExists(c, r.db, project.ID, databaseID)

	if err != nil {
		slog.Error("Error fetching database", err)
		return []models.Rule{}, err
	}

	rules, err := r.db.GetRulesWithEntities(c, database.ID)
	return mapper.ToDomainRulesWithEntity(rules), nil

}

func (r *RuleService) DeleteRule(c context.Context, projectID uuid.UUID, databaseID uuid.UUID, ruleID uuid.UUID) error {
	slog.Info("Deleting rule", "ruleID", ruleID)

	project, err := CheckProjectExists(c, r.db, projectID)

	if err != nil {
		slog.Error("Error fetching project", err)
		return err
	}

	database, err := CheckDatabaseExists(c, r.db, project.ID, databaseID)

	if err != nil {
		slog.Error("Error fetching database", err)
		return err
	}

	err = r.db.DeleteRule(c, sqlc.DeleteRuleParams{
		DatabaseID: database.ID,
		RuleID:     ruleID,
	})
	if err != nil {
		slog.Error("Unable to delete rule", "error", err)
		return err
	}

	return nil
}
