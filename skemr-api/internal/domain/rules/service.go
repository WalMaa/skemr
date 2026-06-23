package rules

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/errormsg"
	"github.com/walmaa/skemr-common/models"
)

type RuleService struct {
	ruleStore     RuleStore
	scopeResolver ScopeResolver
}

type RuleStore interface {
	GetRuleWithEntity(ctx context.Context, params sqlc.GetRuleWithEntityParams) (sqlc.GetRuleWithEntityRow, error)
	GetRuleByDatabaseAndName(ctx context.Context, params sqlc.GetRuleByDatabaseAndNameParams) (sqlc.Rule, error)
	CreateRule(ctx context.Context, dto sqlc.CreateRuleParams) (sqlc.Rule, error)
	GetRulesWithEntities(ctx context.Context, row sqlc.GetRulesWithEntitiesParams) ([]sqlc.GetRulesWithEntitiesRow, error)
	DeleteRule(ctx context.Context, params sqlc.DeleteRuleParams) error
}

type ScopeResolver interface {
	RequireDatabase(c context.Context, projectID uuid.UUID, databaseID uuid.UUID) (models.Database, error)
	RequireDatabaseEntity(c context.Context, projectID uuid.UUID, databaseID uuid.UUID, entityID uuid.UUID) (models.DatabaseEntity, error)
}

func NewRuleService(ruleStore RuleStore, resolver ScopeResolver) *RuleService {
	return &RuleService{ruleStore: ruleStore, scopeResolver: resolver}
}

func (r *RuleService) GetRule(c context.Context, projectID uuid.UUID, databaseID uuid.UUID, ruleID uuid.UUID) (models.Rule, error) {
	slog.Info("Fetching rule", "ruleID", ruleID, "databaseID", databaseID, "projectID", projectID)

	rule, err := r.ruleStore.GetRuleWithEntity(c, sqlc.GetRuleWithEntityParams{
		ProjectID:  projectID,
		DatabaseID: databaseID,
		RuleID:     ruleID,
	})

	if err != nil {
		slog.Error("Unable to fetch rule", "err", err)
		return models.Rule{}, err
	}

	return ToDomainRuleWithEntity(rule), nil
}

func (r *RuleService) CreateRule(c context.Context, projectID uuid.UUID, databaseId uuid.UUID, dto RuleCreationDto) (models.Rule, error) {
	slog.Info("Creating rule", "name", dto.Name, "databaseID", databaseId, "projectID", projectID)

	_, err := r.scopeResolver.RequireDatabase(c, projectID, databaseId)

	if err != nil {
		slog.Error("Error fetching database", "err", err)
		return models.Rule{}, err
	}

	_, err = r.scopeResolver.RequireDatabaseEntity(c, projectID, databaseId, dto.DataBaseEntityId)

	if err != nil {
		slog.Error("Error fetching database entity", "err", err)
		return models.Rule{}, err
	}

	// Check if a rule with the same name already exists
	exists, err := r.ruleStore.GetRuleByDatabaseAndName(c, sqlc.GetRuleByDatabaseAndNameParams{
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

	rule, err := r.ruleStore.CreateRule(c, ToSqlcCreateRule(databaseId, dto))
	if err != nil {
		slog.Error("Unable to create a Rule", "err", err)
		return models.Rule{}, err
	}

	return ToDomainRule(rule), nil
}

func (r *RuleService) ListRulesByDatabase(c context.Context, projectID uuid.UUID, databaseID uuid.UUID) ([]models.Rule, error) {
	slog.Info("Listing rules", "projectID", projectID, "databaseID", databaseID)

	rules, err := r.ruleStore.GetRulesWithEntities(c, sqlc.GetRulesWithEntitiesParams{
		DatabaseID: databaseID,
		ProjectID:  projectID,
	})

	if err != nil {
		slog.Error("Unable to get rules", "err", err)
		return []models.Rule{}, err
	}

	return ToDomainRulesWithEntity(rules), nil
}

func (r *RuleService) DeleteRule(c context.Context, projectID uuid.UUID, databaseId uuid.UUID, ruleID uuid.UUID) error {
	slog.Info("Deleting rule", "ruleID", ruleID)

	_, err := r.scopeResolver.RequireDatabase(c, projectID, databaseId)

	if err != nil {
		slog.Error("Error fetching database", "err", err)
		return err
	}

	err = r.ruleStore.DeleteRule(c, sqlc.DeleteRuleParams{
		DatabaseID: databaseId,
		RuleID:     ruleID,
	})
	if err != nil {
		slog.Error("Unable to delete rule", "error", err)
		return err
	}

	return nil
}
