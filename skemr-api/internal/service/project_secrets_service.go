package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/dto"
	"github.com/walmaa/skemr-api/internal/mapper"
	"github.com/walmaa/skemr-api/internal/tokens"
	"github.com/walmaa/skemr-common/models"
)

type ProjectSecretsService struct {
	db sqlc.Querier
}

func NewProjectSecretsService(q sqlc.Querier) *ProjectSecretsService {
	return &ProjectSecretsService{db: q}
}

func (s *ProjectSecretsService) CreateToken(c context.Context, projectId uuid.UUID, dto dto.SecretCreationDto) (string, error) {
	slog.Info("Creating a secret", "projectId", projectId, "name", dto.Name)

	project, err := CheckProjectExists(c, s.db, projectId)

	if err != nil {
		slog.Error("Unable to get project", err)
		return "", err
	}

	tokenToShow, prefix, secret, err := tokens.GenerateToken(9, 32)

	if err != nil {
		slog.Error("Unable to generate token", err)
		return "", err
	}

	// Generate a hash to store to db. Used for verification
	verifier, err := tokens.HashSecret(secret, tokens.DefaultParams)

	if err != nil {
		slog.Error("Unable to hash token", err)
		return "", err
	}

	// Save the prefix for lookup, and the hash for verification
	_, err = s.db.CreateProjectSecretKey(c, sqlc.CreateProjectSecretKeyParams{
		ProjectID: project.ID,
		Name:      dto.Name,
		Prefix:    prefix,
		Hash:      verifier,
		ExpiresAt: pgtype.Timestamptz{},
	})
	if err != nil {
		slog.Error("Error saving a project access token", err)
		return "", err
	}

	return tokenToShow, nil
}

func (s *ProjectSecretsService) GetTokens(c context.Context, projectId uuid.UUID) ([]models.ProjectAccessToken, error) {
	slog.Info("Getting tokens", "projectId", projectId)
	project, err := CheckProjectExists(c, s.db, projectId)

	if err != nil {
		slog.Error("Unable to get project", err)
		return nil, err
	}

	accessTokens, err := s.db.GetProjectAccessTokens(c, project.ID)
	if err != nil {
		slog.Error("Unable to get tokens", err)
		return nil, err
	}

	return mapper.ToDomainProjectAccessKeys(accessTokens), nil

}
