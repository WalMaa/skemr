package controlplaneclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/walmaa/skemr-common/models"
)

var client = &http.Client{Timeout: 10 * time.Second}

func GetRules(ctx context.Context, projectId string, databaseId string) []models.Rule {
	slog.Info("Fetching rules from control plane", "projectId", projectId, "databaseId", databaseId)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://localhost:8080/api/v1/projects/%s/databases/%s/rules", projectId, databaseId), nil)
	if err != nil {
		slog.Error("Error creating request", err)
		panic(err)
	}
	resp, err := client.Do(req)

	if err != nil {
		slog.Error("HTTP request error", err)
		panic(err)
	}

	defer resp.Body.Close()

	var out []models.Rule
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		slog.Error("Error decoding response body", err)
		panic(err)
	}

	return out
}
