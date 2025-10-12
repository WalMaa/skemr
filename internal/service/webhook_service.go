package service

import (
	"context"
	"log/slog"

	"github.com/walmaa/skemr/db/sqlc"
)

type WebhookService struct {
	db sqlc.Querier
}

func NewWebhookService(q sqlc.Querier) *WebhookService {
	return &WebhookService{db: q}
}

func (s *WebhookService) HandleGitLabWebhook(c context.Context, payload []byte, secret string) error {
	slog.Info("Handling GitLab webhook", "payload", string(payload), "secret", secret)
	// TODO: Map the secret to a project and process the payload accordingly
	project, err := s.db.GetProjectBySecretKey(c, secret)
	if err != nil {
		slog.Warn("Error getting project by secret", "secret", secret, "err", err)
		return err
	}
	slog.Info("Received webhook for project", "project_id", project.ID, "project_name", project.Name)
	return nil
}
