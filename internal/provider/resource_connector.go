package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/iancrichardson/terraform-provider-wiz/internal/client"
)

// Helper function to normalize JSON strings
func normalizeJSON(v interface{}) string {
	if v == nil || v.(string) == "" {
		return ""
	}
	var jsonObj interface{}
	if err := json.Unmarshal([]byte(v.(string)), &jsonObj); err != nil {
		return v.(string)
	}
	jsonBytes, err := json.Marshal(jsonObj)
	if err != nil {
		return v.(string)
	}
	return string(jsonBytes)
}

// Helper function for deep comparison of maps
func deepCompare(current, desired map[string]interface{}, ignoredFields []string) (bool, []string) {
	changedFields := []string{}

	// Check if fields in desired exist in current and have the same value
	for k, desiredVal := range desired {
		// Skip ignored fields
		if containsString(ignoredFields, k) {
			continue
		}

		currentVal, exists := current[k]
		if !exists {
			changedFields = append(changedFields, k)
			continue
		}

		// Handle nested maps
		desiredMap, desiredIsMap := desiredVal.(map[string]interface{})
		currentMap, currentIsMap := currentVal.(map[string]interface{})

		if desiredIsMap && currentIsMap {
			// Recursively compare nested maps
			equal, nestedChanges := deepCompare(currentMap, desiredMap, ignoredFields)
			if !equal {
				for _, nestedField := range nestedChanges {
					changedFields = append(changedFields, fmt.Sprintf("%s.%s", k, nestedField))
				}
			}
			continue
		}

		// Handle JSON strings
		if desiredStr, ok := desiredVal.(string); ok {
			if currentStr, ok := currentVal.(string); ok {
				// Try to parse as JSON and compare
				var desiredJSON, currentJSON interface{}
				desiredErr := json.Unmarshal([]byte(desiredStr), &desiredJSON)
				currentErr := json.Unmarshal([]byte(currentStr), &currentJSON)

				if desiredErr == nil && currentErr == nil {
					// Both are valid JSON, compare as JSON
					desiredBytes, _ := json.Marshal(desiredJSON)
					currentBytes, _ := json.Marshal(currentJSON)
					if string(desiredBytes) != string(currentBytes) {
						changedFields = append(changedFields, k)
					}
					continue
				}
			}
		}

		// Simple equality check for other types
		if !reflect.DeepEqual(currentVal, desiredVal) {
			changedFields = append(changedFields, k)
		}
	}

	return len(changedFields) == 0, changedFields
}

// Helper function to check if a string is in a slice
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// Generate a human-readable diff
func generateDiff(current, desired map[string]interface{}) string {
	var diff strings.Builder
	diff.WriteString("Changes:\n")

	// Check for changes in desired vs current
	for k, desiredVal := range desired {
		currentVal, exists := current[k]
		if !exists {
			diff.WriteString(fmt.Sprintf("+ %s: %v (added)\n", k, desiredVal))
			continue
		}

		if !reflect.DeepEqual(currentVal, desiredVal) {
			diff.WriteString(fmt.Sprintf("~ %s: %v -> %v\n", k, currentVal, desiredVal))
		}
	}

	// Check for fields in current that are not in desired
	for k, currentVal := range current {
		if _, exists := desired[k]; !exists {
			diff.WriteString(fmt.Sprintf("- %s: %v (removed)\n", k, currentVal))
		}
	}

	return diff.String()
}

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
				StateFunc:   normalizeJSON,
			},
			"extra_config": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Extra configuration for the connector in JSON format",
				StateFunc:   normalizeJSON,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the connector",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the connector is enabled",
			},
			"last_activity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp of the last activity for this connector",
			},
			"outpost_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the associated outpost",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

	// Set additional fields
	if status, ok := connector["status"].(string); ok {
		if err := d.Set("status", status); err != nil {
			return diag.FromErr(err)
		}
	}

	if enabled, ok := connector["enabled"].(bool); ok {
		if err := d.Set("enabled", enabled); err != nil {
			return diag.FromErr(err)
		}
	}

	if lastActivity, ok := connector["lastActivity"].(string); ok {
		if err := d.Set("last_activity", lastActivity); err != nil {
			return diag.FromErr(err)
		}
	}

	// Set outpost_id if available
	if outpost, ok := connector["outpost"].(map[string]interface{}); ok && outpost != nil {
		if outpostID, ok := outpost["id"].(string); ok {
			if err := d.Set("outpost_id", outpostID); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return diags
}

func resourceConnectorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	// Get current connector state
	connectorID := d.Id()
	currentConnector, err := c.GetConnector(ctx, connectorID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting current connector state: %w", err))
	}

	// Get desired state from schema
	name := d.Get("name").(string)

	// Parse auth_params and extra_config JSON
	var authParams map[string]interface{}
	if d.HasChange("auth_params") {
		if err := json.Unmarshal([]byte(d.Get("auth_params").(string)), &authParams); err != nil {
			return diag.FromErr(fmt.Errorf("error parsing auth_params: %w", err))
		}
	}

	var extraConfig map[string]interface{}
	if d.HasChange("extra_config") {
		if extraConfigStr, ok := d.Get("extra_config").(string); ok && extraConfigStr != "" {
			if err := json.Unmarshal([]byte(extraConfigStr), &extraConfig); err != nil {
				return diag.FromErr(fmt.Errorf("error parsing extra_config: %w", err))
			}
		}
	}

	// Build desired state for comparison
	desiredState := map[string]interface{}{
		"name": name,
	}

	// Add auth_params and extra_config to desired state if they've changed
	if authParams != nil {
		desiredState["authParams"] = authParams
	} else if currentAuthParams, ok := currentConnector["authParams"]; ok {
		desiredState["authParams"] = currentAuthParams
	}

	if extraConfig != nil {
		desiredState["extraConfig"] = extraConfig
	} else if currentExtraConfig, ok := currentConnector["extraConfig"]; ok {
		desiredState["extraConfig"] = currentExtraConfig
	}

	// Define fields that should be ignored in comparison
	ignoredFields := []string{"id", "status", "lastActivity", "outpost", "type", "config"}

	// Compare current and desired state
	equal, _ := deepCompare(currentConnector, desiredState, ignoredFields)

	// Only update if there are changes
	if !equal {
		// Log the changes
		diff := generateDiff(currentConnector, desiredState)
		fmt.Printf("Updating connector %s with changes: %s\n", connectorID, diff)

		// Update the connector
		if err := c.UpdateConnector(ctx, connectorID, name, authParams, extraConfig); err != nil {
			return diag.FromErr(fmt.Errorf("error updating connector: %w", err))
		}
	} else {
		fmt.Printf("No changes detected for connector %s\n", connectorID)
	}

	return resourceConnectorRead(ctx, d, m)
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
