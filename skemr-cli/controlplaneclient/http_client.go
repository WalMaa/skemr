package controlplaneclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"github.com/walmaa/skemr-common/models"
)

var client = &http.Client{Timeout: 10 * time.Second}

func GetRules(ctx context.Context, projectId string, databaseId string, token string) ([]models.Rule, error) {
	host := viper.GetString("controlPlaneUrl")
	slog.Info("Fetching rules from control plane", "projectId", projectId, "databaseId", databaseId, "host", host)
	var bearer = "Bearer " + token
	var url = fmt.Sprintf("%s/api/v1/projects/%s/databases/%s/rules", host, projectId, databaseId)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		slog.Error("Error creating request", err)
		return nil, err
	}
	req.Header.Add("Authorization", bearer)
	resp, err := client.Do(req)

	if err != nil {
		slog.Error("HTTP request error", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		slog.Error("Error getting rules", "statusCode", strconv.Itoa(resp.StatusCode))
		return nil, fmt.Errorf("Error getting rules, status code: %d", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		slog.Error("Error getting rules", "statusCode", strconv.Itoa(resp.StatusCode))
	}

	defer resp.Body.Close()

	var out []models.Rule
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		slog.Error("Error decoding response body", err)
		return nil, err
	}

	return out, nil
}
