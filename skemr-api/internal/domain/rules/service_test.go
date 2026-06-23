package rules

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/errormsg"
	"github.com/walmaa/skemr-api/test/mocks"
	"github.com/walmaa/skemr-common/models"
)

func TestCreateLockedRule(t *testing.T) {
	projectId := uuid.New()
	databaseId := uuid.New()
	scopeResolver := mocks.NewMockScopeResolver(t)
	ruleStore := mocks.NewMockRuleStore(t)
	ruleType := models.RuleTypeLocked

	svc := NewRuleService(ruleStore, scopeResolver)

	input := RuleCreationDto{
		Name:             "test",
		RuleType:         ruleType,
		Attributes:       models.RuleAttributes{},
		DataBaseEntityId: uuid.UUID{},
	}

	scopeResolver.On("RequireDatabase", mock.Anything, projectId, databaseId).Return(models.Database{}, nil)

	scopeResolver.On("RequireDatabaseEntity", mock.Anything, projectId, databaseId, input.DataBaseEntityId).Return(models.DatabaseEntity{}, nil)
	ruleStore.On("GetRuleByDatabaseAndName", mock.Anything, mock.Anything).Return(sqlc.Rule{}, pgx.ErrNoRows)
	ruleStore.On("CreateRule", mock.Anything, mock.Anything).Return(sqlc.Rule{
		ID:         uuid.New(),
		DatabaseID: databaseId,
		Name:       "test",
		Attributes: nil,
		Type:       sqlc.RuleTypeLocked,
	}, nil)

	rule, err := svc.CreateRule(t.Context(), projectId, databaseId, input)

	require.NoError(t, err)
	require.Equal(t, ruleType, rule.RuleType)
	ruleStore.AssertCalled(t, "CreateRule", mock.Anything, mock.Anything)
}

func TestCreateDeprecatedRule(t *testing.T) {
	projectId := uuid.New()
	databaseId := uuid.New()
	scopeResolver := mocks.NewMockScopeResolver(t)
	ruleStore := mocks.NewMockRuleStore(t)
	ruleType := models.RuleTypeDeprecated

	svc := NewRuleService(ruleStore, scopeResolver)

	input := RuleCreationDto{
		Name:             "test",
		RuleType:         ruleType,
		Attributes:       models.RuleAttributes{},
		DataBaseEntityId: uuid.UUID{},
	}

	scopeResolver.On("RequireDatabase", mock.Anything, projectId, databaseId).Return(models.Database{}, nil)

	scopeResolver.On("RequireDatabaseEntity", mock.Anything, projectId, databaseId, input.DataBaseEntityId).Return(models.DatabaseEntity{}, nil)
	ruleStore.On("GetRuleByDatabaseAndName", mock.Anything, mock.Anything).Return(sqlc.Rule{}, pgx.ErrNoRows)
	ruleStore.On("CreateRule", mock.Anything, mock.Anything).Return(sqlc.Rule{
		ID:         uuid.New(),
		DatabaseID: databaseId,
		Name:       "test",
		Attributes: nil,
		Type:       sqlc.RuleTypeDeprecated,
	}, nil)

	rule, err := svc.CreateRule(t.Context(), projectId, databaseId, input)

	require.NoError(t, err)
	require.Equal(t, ruleType, rule.RuleType)
	ruleStore.AssertCalled(t, "CreateRule", mock.Anything, mock.Anything)
}

func TestCreateAdvisoryRule(t *testing.T) {
	projectId := uuid.New()
	databaseId := uuid.New()
	scopeResolver := mocks.NewMockScopeResolver(t)
	ruleStore := mocks.NewMockRuleStore(t)
	ruleType := models.RuleTypeAdvisory

	svc := NewRuleService(ruleStore, scopeResolver)

	input := RuleCreationDto{
		Name:             "test",
		RuleType:         ruleType,
		Attributes:       models.RuleAttributes{},
		DataBaseEntityId: uuid.UUID{},
	}

	scopeResolver.On("RequireDatabase", mock.Anything, projectId, databaseId).Return(models.Database{}, nil)

	scopeResolver.On("RequireDatabaseEntity", mock.Anything, projectId, databaseId, input.DataBaseEntityId).Return(models.DatabaseEntity{}, nil)
	ruleStore.On("GetRuleByDatabaseAndName", mock.Anything, mock.Anything).Return(sqlc.Rule{}, pgx.ErrNoRows)
	ruleStore.On("CreateRule", mock.Anything, mock.Anything).Return(sqlc.Rule{
		ID:         uuid.New(),
		DatabaseID: databaseId,
		Name:       "test",
		Attributes: nil,
		Type:       sqlc.RuleTypeAdvisory,
	}, nil)

	rule, err := svc.CreateRule(t.Context(), projectId, databaseId, input)

	require.NoError(t, err)
	require.Equal(t, ruleType, rule.RuleType)
	ruleStore.AssertCalled(t, "CreateRule", mock.Anything, mock.Anything)
}

func TestCreateWarningRule(t *testing.T) {
	projectId := uuid.New()
	databaseId := uuid.New()
	scopeResolver := mocks.NewMockScopeResolver(t)
	ruleStore := mocks.NewMockRuleStore(t)
	ruleType := models.RuleTypeWarn

	svc := NewRuleService(ruleStore, scopeResolver)

	input := RuleCreationDto{
		Name:             "test",
		RuleType:         ruleType,
		Attributes:       models.RuleAttributes{},
		DataBaseEntityId: uuid.UUID{},
	}

	scopeResolver.On("RequireDatabase", mock.Anything, projectId, databaseId).Return(models.Database{}, nil)

	scopeResolver.On("RequireDatabaseEntity", mock.Anything, projectId, databaseId, input.DataBaseEntityId).Return(models.DatabaseEntity{}, nil)
	ruleStore.On("GetRuleByDatabaseAndName", mock.Anything, mock.Anything).Return(sqlc.Rule{}, pgx.ErrNoRows)
	ruleStore.On("CreateRule", mock.Anything, mock.Anything).Return(sqlc.Rule{
		ID:         uuid.New(),
		DatabaseID: databaseId,
		Name:       "test",
		Attributes: nil,
		Type:       sqlc.RuleTypeWarn,
	}, nil)

	rule, err := svc.CreateRule(t.Context(), projectId, databaseId, input)

	require.NoError(t, err)
	require.Equal(t, ruleType, rule.RuleType)
	ruleStore.AssertCalled(t, "CreateRule", mock.Anything, mock.Anything)
}

func TestCreateRuleWithoutDatabaseEntity(t *testing.T) {
	projectId := uuid.New()
	databaseId := uuid.New()
	scopeResolver := mocks.NewMockScopeResolver(t)
	ruleStore := mocks.NewMockRuleStore(t)

	svc := NewRuleService(ruleStore, scopeResolver)

	input := RuleCreationDto{
		Name:             "test",
		RuleType:         models.RuleTypeLocked,
		Attributes:       models.RuleAttributes{},
		DataBaseEntityId: uuid.UUID{},
	}

	scopeResolver.On("RequireDatabase", mock.Anything, projectId, databaseId).Return(models.Database{}, nil)

	scopeResolver.On("RequireDatabaseEntity", mock.Anything, projectId, databaseId, input.DataBaseEntityId).Return(models.DatabaseEntity{}, &models.ErrorResponse{
		Message: errormsg.ErrDatabaseEntityNotFound,
		Errors:  nil,
		Status:  http.StatusBadRequest,
	})

	_, err := svc.CreateRule(t.Context(), projectId, databaseId, input)
	require.Error(t, err)
	require.Equal(t, errormsg.ErrDatabaseEntityNotFound, err.Error())

	ruleStore.AssertNotCalled(t, "CreateRule", mock.Anything, mock.Anything)
}

func TestCreateRuleWithNonExistentDatabaseEntity(t *testing.T) {
	projectId := uuid.New()
	databaseId := uuid.New()
	scopeResolver := mocks.NewMockScopeResolver(t)
	ruleStore := mocks.NewMockRuleStore(t)

	svc := NewRuleService(ruleStore, scopeResolver)

	input := RuleCreationDto{
		Name:             "test",
		RuleType:         models.RuleTypeLocked,
		Attributes:       models.RuleAttributes{},
		DataBaseEntityId: uuid.New(),
	}

	scopeResolver.On("RequireDatabase", mock.Anything, projectId, databaseId).Return(models.Database{}, nil)

	scopeResolver.On("RequireDatabaseEntity", mock.Anything, projectId, databaseId, input.DataBaseEntityId).Return(models.DatabaseEntity{}, &models.ErrorResponse{
		Message: errormsg.ErrDatabaseEntityNotFound,
		Errors:  nil,
		Status:  http.StatusBadRequest,
	})

	_, err := svc.CreateRule(t.Context(), projectId, databaseId, input)
	require.Error(t, err)
	require.Equal(t, errormsg.ErrDatabaseEntityNotFound, err.Error())

	ruleStore.AssertNotCalled(t, "CreateRule", mock.Anything, mock.Anything)

}

func TestCreateRuleWithNonExistentProjectOrDatabase(t *testing.T) {
	projectId := uuid.New()
	databaseId := uuid.New()
	scopeResolver := mocks.NewMockScopeResolver(t)
	ruleStore := mocks.NewMockRuleStore(t)

	svc := NewRuleService(ruleStore, scopeResolver)

	input := RuleCreationDto{
		Name:             "test",
		RuleType:         models.RuleTypeLocked,
		Attributes:       models.RuleAttributes{},
		DataBaseEntityId: uuid.New(),
	}

	scopeResolver.On("RequireDatabase", mock.Anything, projectId, databaseId).Return(models.Database{}, &models.ErrorResponse{
		Message: errormsg.ErrDatabaseNotFound,
		Errors:  nil,
		Status:  http.StatusBadRequest,
	})

	_, err := svc.CreateRule(t.Context(), projectId, databaseId, input)
	require.Error(t, err)
	require.Equal(t, errormsg.ErrDatabaseNotFound, err.Error())

	ruleStore.AssertNotCalled(t, "CreateRule", mock.Anything, mock.Anything)
	ruleStore.AssertNotCalled(t, "RequireDatabaseEntity", mock.Anything, projectId, databaseId, input.DataBaseEntityId)
}

func TestCreateRuleWithExistingName(t *testing.T) {
	projectId := uuid.New()
	databaseId := uuid.New()
	scopeResolver := mocks.NewMockScopeResolver(t)
	ruleStore := mocks.NewMockRuleStore(t)

	svc := NewRuleService(ruleStore, scopeResolver)

	input := RuleCreationDto{
		Name:             "test",
		RuleType:         models.RuleTypeLocked,
		Attributes:       models.RuleAttributes{},
		DataBaseEntityId: uuid.New(),
	}

	scopeResolver.On("RequireDatabase", mock.Anything, projectId, databaseId).Return(models.Database{}, nil)

	scopeResolver.On("RequireDatabaseEntity", mock.Anything, projectId, databaseId, input.DataBaseEntityId).Return(models.DatabaseEntity{}, nil)

	ruleStore.On("GetRuleByDatabaseAndName", mock.Anything, sqlc.GetRuleByDatabaseAndNameParams{
		DatabaseID: databaseId,
		Name:       "test",
	}).Return(sqlc.Rule{
		ID:         uuid.New(),
		DatabaseID: databaseId,
		Name:       "test",
	}, nil)

	_, err := svc.CreateRule(t.Context(), projectId, databaseId, input)
	require.Error(t, err)
	require.Equal(t, errormsg.ErrRuleWithSameName, err.Error())

	ruleStore.AssertNotCalled(t, "CreateRule", mock.Anything, mock.Anything)
}
