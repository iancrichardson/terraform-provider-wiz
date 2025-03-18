package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/iancrichardson/terraform-provider-wiz/internal/client"
)

func dataSourceConnectorConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectorConfigRead,
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of the connector (e.g., azure, aws, gcp)",
			},
			"auth_params": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Authentication parameters for the connector in JSON format",
				StateFunc: func(v interface{}) string {
					// Normalize the JSON string
					var jsonObj interface{}
					json.Unmarshal([]byte(v.(string)), &jsonObj)
					jsonBytes, _ := json.Marshal(jsonObj)
					return string(jsonBytes)
				},
			},
			"extra_config": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Extra configuration for the connector in JSON format",
				StateFunc: func(v interface{}) string {
					if v == nil || v.(string) == "" {
						return ""
					}
					// Normalize the JSON string
					var jsonObj interface{}
					json.Unmarshal([]byte(v.(string)), &jsonObj)
					jsonBytes, _ := json.Marshal(jsonObj)
					return string(jsonBytes)
				},
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of an existing connector to test",
			},
			"success": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the connector configuration is valid",
			},
		},
	}
}

func dataSourceConnectorConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics

	connectorType := d.Get("type").(string)

	// Parse auth_params JSON
	var authParams map[string]interface{}
	if err := json.Unmarshal([]byte(d.Get("auth_params").(string)), &authParams); err != nil {
		return diag.FromErr(fmt.Errorf("error parsing auth_params: %w", err))
	}

	// Parse extra_config JSON if provided
	var extraConfig map[string]interface{}
	if extraConfigStr, ok := d.Get("extra_config").(string); ok && extraConfigStr != "" {
		if err := json.Unmarshal([]byte(extraConfigStr), &extraConfig); err != nil {
			return diag.FromErr(fmt.Errorf("error parsing extra_config: %w", err))
		}
	}

	// Get connector ID if provided
	var id string
	if v, ok := d.Get("id").(string); ok {
		id = v
	}

	// Test the connector configuration
	success, err := c.TestConnectorConfig(ctx, connectorType, authParams, extraConfig, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error testing connector configuration: %w", err))
	}

	// Set the computed values
	if err := d.Set("success", success); err != nil {
		return diag.FromErr(err)
	}

	// Generate a unique ID for the data source
	d.SetId(fmt.Sprintf("%s-%s", connectorType, id))

	return diags
}
