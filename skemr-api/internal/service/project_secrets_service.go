package service

import "github.com/walmaa/skemr-api/db/sqlc"

type ProjectSecretsService struct {
	db sqlc.Querier
}

func NewProjectSecretsService(q sqlc.Querier) *ProjectSecretsService {
	return &ProjectSecretsService{db: q}
}
