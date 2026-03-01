package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/dto"
	"github.com/walmaa/skemr-api/internal/errormsg"
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
	expires := pgtype.Timestamptz{
		Time:  time.Time{},
		Valid: false,
	}
	if dto.ExpiresAt != "" {
		expiry, err := time.Parse(time.RFC3339, dto.ExpiresAt)
		if err != nil {
			slog.Error("Unable to parse expiry time", err)
			return "", err
		}
		if expiry.Before(time.Now()) {
			slog.Error("Expiry time is in the past")
			err = fmt.Errorf("expiry time is in the past")
			return "", &errormsg.ErrorResponse{
				Message: errormsg.ErrExpiryTimeInPast.Error(),
				Status:  http.StatusBadRequest,
			}
		}
		expires = pgtype.Timestamptz{
			Time:  expiry,
			Valid: true,
		}

	}

	// Save the prefix for lookup, and the hash for verification
	_, err = s.db.CreateProjectSecretKey(c, sqlc.CreateProjectSecretKeyParams{
		ProjectID: project.ID,
		Name:      dto.Name,
		Prefix:    prefix,
		Hash:      verifier,
		ExpiresAt: expires,
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

func (s *ProjectSecretsService) DeleteToken(c context.Context, projectId uuid.UUID, secretId uuid.UUID) error {
	slog.Info("Deleting token", "projectId", projectId, "secretId", secretId)

	project, err := CheckProjectExists(c, s.db, projectId)

	if err != nil {
		slog.Error("Unable to get project")
		return err
	}

	err = s.db.DeleteProjectAccessToken(c, sqlc.DeleteProjectAccessTokenParams{
		ProjectID: project.ID,
		SecretID:  secretId,
	})
	if err != nil {
		slog.Error("Unable to delete project access token", err)
		return err
	}

	return nil

}
