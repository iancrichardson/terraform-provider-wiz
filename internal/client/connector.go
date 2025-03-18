package client

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// TestConnectorConfigResponse represents the response from the testConnectorConfig mutation
type TestConnectorConfigResponse struct {
	TestConnectorConfig struct {
		Success bool `json:"success"`
	} `json:"testConnectorConfig"`
}

// CreateConnectorResponse represents the response from the createConnector mutation
type CreateConnectorResponse struct {
	CreateConnector struct {
		Connector struct {
			ID         string      `json:"id"`
			Name       string      `json:"name"`
			AuthParams interface{} `json:"authParams"`
			Type       struct {
				ID string `json:"id"`
			} `json:"type"`
			ExtraConfig interface{} `json:"extraConfig"`
			Outpost     struct {
				ID             string `json:"id"`
				ServiceAccount struct {
					ClientID     string `json:"clientId"`
					ClientSecret string `json:"clientSecret"`
				} `json:"serviceAccount"`
			} `json:"outpost"`
		} `json:"connector"`
	} `json:"createConnector"`
}

// DeleteConnectorResponse represents the response from the deleteConnector mutation
type DeleteConnectorResponse struct {
	DeleteConnector struct {
		Stub string `json:"_stub"`
	} `json:"deleteConnector"`
}

// TestConnectorConfig tests a connector configuration
func (c *Client) TestConnectorConfig(ctx context.Context, connectorType string, authParams map[string]interface{}, extraConfig map[string]interface{}, id string) (bool, error) {
	query := `
		query TestConnectorConfig($type: ID!, $authParams: JSON!, $extraConfig: JSON, $id: String) {
			testConnectorConfig(
				type: $type
				authParams: $authParams
				extraConfig: $extraConfig
				id: $id
			) {
				success
			}
		}
	`

	variables := map[string]interface{}{
		"type":       connectorType,
		"authParams": authParams,
	}

	if extraConfig != nil {
		variables["extraConfig"] = extraConfig
	}

	if id != "" {
		variables["id"] = id
	}

	var response TestConnectorConfigResponse
	if err := c.RunQuery(ctx, query, variables, &response); err != nil {
		return false, fmt.Errorf("error testing connector config: %w", err)
	}

	return response.TestConnectorConfig.Success, nil
}

// CreateConnector creates a new connector
func (c *Client) CreateConnector(ctx context.Context, name string, connectorType string, authParams map[string]interface{}, extraConfig map[string]interface{}) (string, error) {
	query := `
		mutation CreateConnector($input: CreateConnectorInput!) {
			createConnector(input: $input) {
				connector {
					id
					name
					authParams
					type {
						id
					}
					extraConfig
					outpost {
						id
						serviceAccount {
							clientId
							clientSecret
						}
					}
				}
			}
		}
	`

	input := map[string]interface{}{
		"name":       name,
		"type":       connectorType,
		"authParams": authParams,
	}

	if extraConfig != nil {
		input["extraConfig"] = extraConfig
	}

	variables := map[string]interface{}{
		"input": input,
	}

	var response CreateConnectorResponse
	if err := c.RunQuery(ctx, query, variables, &response); err != nil {
		return "", fmt.Errorf("error creating connector: %w", err)
	}

	return response.CreateConnector.Connector.ID, nil
}

// DeleteConnector deletes a connector
func (c *Client) DeleteConnector(ctx context.Context, id string) error {
	query := `
		mutation DeleteConnector($input: DeleteConnectorInput!) {
			deleteConnector(input: $input) {
				_stub
			}
		}
	`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"id": id,
		},
	}

	var response DeleteConnectorResponse
	if err := c.RunQuery(ctx, query, variables, &response); err != nil {
		return fmt.Errorf("error deleting connector: %w", err)
	}

	return nil
}

// GetConnectorResponse represents the response from the getConnector query
type GetConnectorResponse struct {
	Connector struct {
		ID           string      `json:"id"`
		Name         string      `json:"name"`
		AuthParams   interface{} `json:"authParams"`
		ExtraConfig  interface{} `json:"extraConfig"`
		Status       string      `json:"status"`
		Enabled      bool        `json:"enabled"`
		LastActivity string      `json:"lastActivity"`
		Outpost      struct {
			ID     string `json:"id"`
			Config struct {
				Environment string `json:"environment"`
			} `json:"config"`
		} `json:"outpost"`
		Config interface{} `json:"config"`
		Type   struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"type"`
	} `json:"connector"`
}

// UpdateConnectorResponse represents the response from the updateConnector mutation
type UpdateConnectorResponse struct {
	UpdateConnector struct {
		Connector struct {
			ID           string      `json:"id"`
			Name         string      `json:"name"`
			Status       string      `json:"status"`
			Enabled      bool        `json:"enabled"`
			LastActivity string      `json:"lastActivity"`
			ExtraConfig  interface{} `json:"extraConfig"`
		} `json:"connector"`
	} `json:"updateConnector"`
}

// GetConnector gets a connector by ID with detailed information
func (c *Client) GetConnector(ctx context.Context, id string) (map[string]interface{}, error) {
	query := `
		query GetConnector($connectorId: ID!) {
		  connector(id: $connectorId) {
			id
			name
			status
			enabled
			lastActivity
			authParams
			extraConfig
			outpost {
			  id
			  config {
				... on OutpostAzureConfig {
				  environment
				}
			  }
			}
			config {
			  ... on ConnectorConfigAWS {
				region
				customerRoleARN
				scheduledSecurityToolScanningSettings {
				  enabled
				  publicBucketsScanningEnabled
				}
			  }
			  ... on ConnectorConfigGCP {
				isManagedIdentity
				projects
				excludedProjects
				includedFolders
				excludedFolders
				organizationId: organization_id
				projectId: project_id
				folderId: folder_id
				customerId: customer_id
				auditLogMonitorEnabled
				scheduledSecurityToolScanningSettings {
				  enabled
				  publicBucketsScanningEnabled
				}
				auditLogsConfig {
				  pub_sub {
					topicName
					subscriptionID
				  }
				}
			  }
			  ... on ConnectorConfigAzure {
				monitorEventHubConnectionString
				excludedSubscriptions
				includedSubscriptions
				excludedManagementGroups
				includedManagementGroups
				auditLogMonitorEnabled
				snapshotsResourceGroupId
				environment
				scheduledSecurityToolScanningSettings {
				  enabled
				  publicBucketsScanningEnabled
				}
				tenantId
				groupId
				subscriptionId
				isManagedIdentity
				isAzureActiveDirectoryOnly
				azureMonitorConfig {
				  eventHub {
					connectionMethod
					name
					namespace
					namespaceTag
				  }
				}
				costAndUsageReportConfig {
				  subscription
				  areStorageSettingsShared
				  amortizedReportConfig {
					exportResourceGroup
					exportStorageAccountName
					exportContainer
					exportDirectory
					exportName
				  }
				  actualReportConfig {
					exportResourceGroup
					exportStorageAccountName
					exportContainer
					exportDirectory
					exportName
				  }
				  isEnabled
				}
			  }
			}
			type {
			  id
			  name
			}
		  }
		}
	`

	variables := map[string]interface{}{
		"connectorId": id,
	}

	var response map[string]interface{}
	err := retryWithBackoff(ctx, func() error {
		return c.RunQuery(ctx, query, variables, &response)
	})

	if err != nil {
		return nil, fmt.Errorf("error getting connector: %w", err)
	}

	// Check if connector exists
	connectorData, ok := response["connector"]
	if !ok || connectorData == nil {
		return nil, fmt.Errorf("connector not found: %s", id)
	}

	// Convert to map
	connectorMap, ok := connectorData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid connector data format")
	}

	return connectorMap, nil
}

// UpdateConnector updates an existing connector
func (c *Client) UpdateConnector(ctx context.Context, id string, name string, authParams map[string]interface{}, extraConfig map[string]interface{}) error {
	query := `
		mutation UpdateConnector($input: UpdateConnectorInput!) {
		  updateConnector(input: $input) {
			connector {
			  id
			  name
			  status
			  enabled
			  lastActivity
			  extraConfig
			}
		  }
		}
	`

	// Build the patch object with the changes
	patch := map[string]interface{}{
		"name": name,
	}

	if authParams != nil {
		patch["authParams"] = authParams
	}

	if extraConfig != nil {
		patch["extraConfig"] = extraConfig
	}

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"id":    id,
			"patch": patch,
		},
	}

	var response UpdateConnectorResponse
	err := retryWithBackoff(ctx, func() error {
		return c.RunQuery(ctx, query, variables, &response)
	})

	if err != nil {
		return fmt.Errorf("error updating connector: %w", err)
	}

	return nil
}

// retryWithBackoff retries a function with exponential backoff
func retryWithBackoff(ctx context.Context, f func() error) error {
	var err error
	maxRetries := 5
	baseDelay := 1 * time.Second

	for i := 0; i < maxRetries; i++ {
		err = f()
		if err == nil {
			return nil
		}

		// Check if error is retryable
		if !isRetryableError(err) {
			return err
		}

		// Calculate backoff with jitter
		delay := baseDelay * time.Duration(1<<uint(i)) // Exponential
		jitter := time.Duration(rand.Int63n(int64(delay) / 2))
		delay = delay + jitter

		select {
		case <-time.After(delay):
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("maximum retries exceeded: %w", err)
}

// isRetryableError determines if an error is retryable
func isRetryableError(err error) bool {
	// Check for network errors, rate limits, and temporary API issues
	if strings.Contains(err.Error(), "rate limit") {
		return true
	}
	if strings.Contains(err.Error(), "timeout") {
		return true
	}
	if strings.Contains(err.Error(), "connection reset") {
		return true
	}
	if strings.Contains(err.Error(), "temporary") {
		return true
	}
	if strings.Contains(err.Error(), "service unavailable") {
		return true
	}
	return false
}
