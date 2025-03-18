terraform {
  required_providers {
    wiz = {
      source = "iancrichardson/wiz"
      version = "0.1.0"
    }
  }
}

provider "wiz" {
  client_id     = "SERVICE_ACCOUNT_CLIENT_ID"
  client_secret = "SERVICE_ACCOUNT_CLIENT_SECRET"
  # Optional: Uncomment to override default endpoints
  # api_url       = "https://api.eu1.demo.wiz.io/graphql"
  # auth_url      = "https://auth.demo.wiz.io/oauth/token"
}

# Test a connector configuration
data "wiz_connector_config" "test" {
  type = "azure"
  auth_params = jsonencode({
    isManagedIdentity = true
    subscriptionId    = "2068e3a6-f96a-45d1-aea7-442cd9b2c26d"
    tenantId          = "c76e6a8f-b1ba-44c4-a73f-8928b943f202"
    environment       = "AzurePublicCloud"
  })
}

# Create a connector
resource "wiz_connector" "azure" {
  name = "Azure Connector"
  type = "azure"
  
  auth_params = jsonencode({
    isManagedIdentity = true
    subscriptionId    = "2068e3a6-f96a-45d1-aea7-442cd9b2c26d"
    tenantId          = "c76e6a8f-b1ba-44c4-a73f-8928b943f202"
    environment       = "AzurePublicCloud"
  })
  
  extra_config = jsonencode({
    includedSubscriptions = []
    excludedSubscriptions = []
    includedManagementGroups = []
    excludedManagementGroups = []
    snapshotsResourceGroupId = ""
    auditLogMonitorEnabled = false
    scheduledSecurityToolScanningSettings = {
      enabled = true
      publicBucketsScanningEnabled = false
    }
  })
}

# Output the connector ID
output "connector_id" {
  value = wiz_connector.azure.id
}

# Output the configuration test result
output "config_test_success" {
  value = data.wiz_connector_config.test.success
}
