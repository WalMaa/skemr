package accesstoken

import (
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/mapper"
	"github.com/walmaa/skemr-common/models"
)

func ToDomainProjectAccessKey(token sqlc.ProjectAccessToken) models.ProjectAccessToken {
	return models.ProjectAccessToken{
		ID:        token.ID,
		ProjectID: token.ProjectID,
		Name:      token.Name,
		LastUsed:  mapper.TimePtr(&token.LastUsed),
		ExpiresAt: mapper.TimePtr(&token.ExpiresAt),
		CreatedAt: mapper.Time(&token.CreatedAt),
		UpdatedAt: mapper.Time(&token.UpdatedAt),
	}
}

func ToDomainProjectAccessKeys(s []sqlc.ProjectAccessToken) []models.ProjectAccessToken {
	tokens := make([]models.ProjectAccessToken, len(s))

	for i, token := range s {
		tokens[i] = ToDomainProjectAccessKey(token)
	}
	return tokens
}
