package controlplaneclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"github.com/walmaa/skemr-common/models"
)

var client = &http.Client{Timeout: 10 * time.Second}

func GetRules(ctx context.Context, projectId string, databaseId string, token string) ([]models.Rule, error) {
	host := viper.GetString("controlPlaneUrl")
	slog.Info("Fetching rules from control plane", "projectId", projectId, "databaseId", databaseId, "host", host)
	bearer := "Bearer " + token
	url := fmt.Sprintf("%s/api/v1/projects/%s/databases/%s/integrations/ci-cd/rules", host, projectId, databaseId)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}
	req.Header.Add("Authorization", bearer)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Error closing response body", "error", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var errorResponse models.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, fmt.Errorf("Error getting rules, status code: %d, and error decoding response body: %v", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("Error getting rules, status code: %d, error: %s", resp.StatusCode, errorResponse.Error())
	}

	var out []models.Rule
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	return out, nil
}

func GetDatabaseEntity(ctx context.Context, projectId string, databaseId string, databaseEntityId string, token string) (*models.DatabaseEntity, error) {
	host := viper.GetString("controlPlaneUrl")
	slog.Info("Fetching database entity from control plane", "projectId", projectId, "databaseId", databaseId, "databaseEntityId", databaseEntityId, "host", host)
	bearer := "Bearer " + token
	url := fmt.Sprintf("%s/api/v1/projects/%s/databases/%s/entities/%s", host, projectId, databaseId, databaseEntityId)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}
	req.Header.Add("Authorization", bearer)
	resp, err := client.Do(req)

	if err != nil {
		slog.Error("HTTP request error", "error", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting database entity, status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	var out models.DatabaseEntity
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		slog.Error("Error decoding response body", "error", err)
		return nil, err
	}

	return &out, nil
}

func GetDatabaseEntities(ctx context.Context, projectId string, databaseId string, token string) ([]models.DatabaseEntity, error) {
	host := viper.GetString("controlPlaneUrl")
	slog.Info("Fetching database entities from control plane", "projectId", projectId, "databaseId", databaseId, "host", host)
	bearer := "Bearer " + token
	url := fmt.Sprintf("%s/api/v1/projects/%s/databases/%s/entities", host, projectId, databaseId)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", bearer)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting database entities, status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	var out []models.DatabaseEntity
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	return out, nil
}
