package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/walmaa/skemr-api/db/sqlc"
	"gitlab.com/gitlab-org/api/client-go"
)

type WebhookService struct {
	db sqlc.Querier
}

func NewWebhookService(q sqlc.Querier) *WebhookService {
	return &WebhookService{db: q}
}

func (s *WebhookService) HandleGitLabWebhook(c context.Context, r *http.Request) error {
	slog.Info("Handling GitLab webhook")
	webhookSecret := r.Header.Get("X-Gitlab-Token")

	if webhookSecret == "" {
		slog.Warn("Missing X-Gitlab-Token header")
		return fmt.Errorf("missing X-Gitlab-Token header")
	}

	// TODO: Map the secret to a project and process the payload accordingly
	project, err := s.db.GetProjectBySecretPrefix(c, webhookSecret)
	if err != nil {
		slog.Warn("Error getting project by secret", "err", err)
		//return err
	}
	eventType := gitlab.WebhookEventType(r)
	slog.Info("Got webhook event", "eventType", eventType)
	body, err := io.ReadAll(r.Body)

	if err != nil {
		slog.Warn("Error reading body", "err", err)
		return err
	}

	event, err := gitlab.ParseWebhook(eventType, body)

	if err != nil {
		slog.Warn("Error parsing webhook", "err", err)
		return err
	}

	switch eventType {
	case gitlab.EventTypeMergeRequest:
		slog.Info("Handling merge request for project", "project_name", project.Name)
		if mergeEvent, ok := event.(*gitlab.MergeEvent); ok {
			s.handleGitLabMergeRequest(project, mergeEvent)
		} else {
			slog.Error("Unexpected event", "event", event)
			return fmt.Errorf("unexpected event type %s", eventType)
		}
	default:
		slog.Warn("Unsupported webhook event", "event", eventType)

	}

	return nil
}

func (s *WebhookService) handleGitLabMergeRequest(project sqlc.Project, event *gitlab.MergeEvent) error {
	slog.Info("Handling merge request for project", "project_name", project.Name)
	return nil
}
