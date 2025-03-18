package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/iancrichardson/terraform-provider-wiz/internal/client"
)

func resourceConnector() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectorCreate,
		ReadContext:   resourceConnectorRead,
		UpdateContext: resourceConnectorUpdate,
		DeleteContext: resourceConnectorDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the connector",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceConnectorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics

	name := d.Get("name").(string)
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

	// Test the connector configuration first
	success, err := c.TestConnectorConfig(ctx, connectorType, authParams, extraConfig, "")
	if err != nil {
		return diag.FromErr(fmt.Errorf("error testing connector configuration: %w", err))
	}

	if !success {
		return diag.FromErr(fmt.Errorf("connector configuration test failed"))
	}

	// Create the connector
	id, err := c.CreateConnector(ctx, name, connectorType, authParams, extraConfig)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating connector: %w", err))
	}

	d.SetId(id)

	return diags
}

func resourceConnectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics

	connectorID := d.Id()

	connector, err := c.GetConnector(ctx, connectorID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting connector: %w", err))
	}

	// Set the resource data
	if err := d.Set("name", connector["name"]); err != nil {
		return diag.FromErr(err)
	}

	if typeData, ok := connector["type"].(map[string]interface{}); ok {
		if err := d.Set("type", typeData["id"]); err != nil {
			return diag.FromErr(err)
		}
	}

	// Convert auth_params to JSON string
	if authParams, ok := connector["authParams"]; ok && authParams != nil {
		authParamsJSON, err := json.Marshal(authParams)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error marshaling auth_params: %w", err))
		}
		if err := d.Set("auth_params", string(authParamsJSON)); err != nil {
			return diag.FromErr(err)
		}
	}

	// Convert extra_config to JSON string
	if extraConfig, ok := connector["extraConfig"]; ok && extraConfig != nil {
		extraConfigJSON, err := json.Marshal(extraConfig)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error marshaling extra_config: %w", err))
		}
		if err := d.Set("extra_config", string(extraConfigJSON)); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceConnectorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Updates are not supported by the Wiz API for connectors
	// If any changes are detected, we'll recreate the resource
	return resourceConnectorCreate(ctx, d, m)
}

func resourceConnectorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics

	connectorID := d.Id()

	if err := c.DeleteConnector(ctx, connectorID); err != nil {
		return diag.FromErr(fmt.Errorf("error deleting connector: %w", err))
	}

	d.SetId("")

	return diags
}
