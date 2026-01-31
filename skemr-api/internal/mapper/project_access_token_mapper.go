package mapper

import (
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-common/models"
)

func ToDomainProjectAccessKey(token sqlc.ProjectAccessToken) models.ProjectAccessToken {
	return models.ProjectAccessToken{
		ID:        token.ID,
		ProjectID: token.ProjectID,
		Name:      token.Name,
		LastUsed:  Time(&token.LastUsed),
		ExpiresAt: Time(&token.ExpiresAt),
		CreatedAt: Time(&token.CreatedAt),
		UpdatedAt: Time(&token.UpdatedAt),
	}
}

func ToDomainProjectAccessKeys(s []sqlc.ProjectAccessToken) []models.ProjectAccessToken {
	tokens := make([]models.ProjectAccessToken, len(s))

	for i, token := range s {
		tokens[i] = ToDomainProjectAccessKey(token)
	}
	return tokens
}
