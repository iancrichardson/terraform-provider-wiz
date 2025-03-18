package client

import (
	"context"
	"fmt"
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

// GetConnector gets a connector by ID
func (c *Client) GetConnector(ctx context.Context, id string) (map[string]interface{}, error) {
	query := `
		query GetConnector($id: ID!) {
			connector(id: $id) {
				id
				name
				authParams
				type {
					id
				}
				extraConfig
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response map[string]interface{}
	if err := c.RunQuery(ctx, query, variables, &response); err != nil {
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
