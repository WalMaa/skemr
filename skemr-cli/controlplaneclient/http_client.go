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
	bearer := "Bearer " + token
	url := fmt.Sprintf("%s/api/v1/projects/%s/databases/%s/rules", host, projectId, databaseId)
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
		slog.Error("Error getting rules", "statusCode", strconv.Itoa(resp.StatusCode))
		return nil, fmt.Errorf("Error getting rules, status code: %d", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		slog.Error("Error getting rules", "statusCode", strconv.Itoa(resp.StatusCode))
	}

	defer resp.Body.Close()

	var out []models.Rule
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		slog.Error("Error decoding response body", "error", err)
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
		slog.Error("Error getting database entity", "statusCode", strconv.Itoa(resp.StatusCode))
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
		slog.Error("Error getting database entities", "statusCode", strconv.Itoa(resp.StatusCode))
		return nil, fmt.Errorf("Error getting database entities, status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	var out []models.DatabaseEntity
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		slog.Error("Error decoding response body", "error", err)
		return nil, err
	}

	return out, nil
}
