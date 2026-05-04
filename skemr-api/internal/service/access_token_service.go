package service

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/dto"
	"github.com/walmaa/skemr-api/internal/errormsg"
	"github.com/walmaa/skemr-api/internal/mapper"
	"github.com/walmaa/skemr-api/internal/tokens"
	"github.com/walmaa/skemr-common/models"
)

type AccessTokenService struct {
	db sqlc.Querier
}

const prefixLength = 9
const secretLength = 32

func NewAccessTokenService(q sqlc.Querier) *AccessTokenService {
	return &AccessTokenService{db: q}
}

func (s *AccessTokenService) CreateToken(c context.Context, projectId uuid.UUID, dto dto.SecretCreationDto) (string, error) {
	slog.Info("Creating a secret", "projectId", projectId, "name", dto.Name)

	project, err := CheckProjectExists(c, s.db, projectId)

	if err != nil {
		slog.Error("Unable to get project", "err", err)
		return "", err
	}

	tokenToShow, prefix, secret, err := tokens.GenerateToken(prefixLength, secretLength)

	if err != nil {
		slog.Error("Unable to generate token", "err", err)
		return "", err
	}

	// Generate a hash to store to db. Used for verification
	verifier, err := tokens.HashSecret(secret, tokens.DefaultParams)

	if err != nil {
		slog.Error("Unable to hash token", "err", err)
		return "", err
	}
	expires := pgtype.Timestamptz{
		Time:  time.Time{},
		Valid: false,
	}
	if dto.ExpiresAt != "" {
		expiry, err := time.Parse(time.RFC3339, dto.ExpiresAt)
		if err != nil {
			slog.Error("Unable to parse expiry time", "err", err)
			return "", err
		}
		if expiry.Before(time.Now()) {
			slog.Error("Expiry time is in the past")
			return "", &models.ErrorResponse{
				Message: errormsg.ErrExpiryTimeInPast,
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
		slog.Error("Error saving a project access token", "err", err)
		return "", err
	}

	return tokenToShow, nil
}

func (s *AccessTokenService) GetTokens(c context.Context, projectId uuid.UUID) ([]models.ProjectAccessToken, error) {
	slog.Info("Getting tokens", "projectId", projectId)
	project, err := CheckProjectExists(c, s.db, projectId)

	if err != nil {
		slog.Error("Unable to get project", "err", err)
		return nil, err
	}

	accessTokens, err := s.db.GetProjectAccessTokens(c, project.ID)
	if err != nil {
		slog.Error("Unable to get tokens", "err", err)
		return nil, err
	}

	return mapper.ToDomainProjectAccessKeys(accessTokens), nil

}

func (s *AccessTokenService) DeleteToken(c context.Context, projectId uuid.UUID, secretId uuid.UUID) error {
	slog.Info("Deleting token", "projectId", projectId, "secretId", secretId)

	project, err := CheckProjectExists(c, s.db, projectId)

	if err != nil {
		slog.Error("Unable to get project", "err", err)
		return err
	}

	err = s.db.DeleteProjectAccessToken(c, sqlc.DeleteProjectAccessTokenParams{
		ProjectID: project.ID,
		SecretID:  secretId,
	})
	if err != nil {
		slog.Error("Unable to delete project access token", "err", err)
		return err
	}

	return nil

}

func (s *AccessTokenService) ValidateToken(c context.Context, projectId uuid.UUID, token string) (bool, error) {
	slog.Info("Validating token")

	// Extract the prefix to find the token in the database

	if len(token) < prefixLength {
		slog.Error("Token is too short to be valid")
		return false, nil
	}

	prefix := strings.Split(token, ".")[0]
	secretPart := strings.Split(token, ".")[1]

	if prefix == "" {
		slog.Error("Token does not contain a prefix")
		return false, nil
	}

	if secretPart == "" {
		slog.Error("Token does not contain a secret part")
		return false, nil
	}

	hash, err := s.db.GetHashByPrefixAndProjectID(c, sqlc.GetHashByPrefixAndProjectIDParams{
		Prefix:    prefix,
		ProjectID: projectId,
	})

	if errors.Is(err, pgx.ErrNoRows) {
		slog.Info("No token found with the given prefix")
		return false, nil
	}

	if err != nil {
		slog.Error("Unable to get token hash from database", "err", err)
		return false, err
	}

	ok, err := tokens.VerifySecret(secretPart, hash)

	if err != nil {
		slog.Error("Error verifying token", "err", err)
		return false, err
	}

	return ok, nil
}
