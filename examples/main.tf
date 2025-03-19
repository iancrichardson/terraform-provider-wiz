terraform {
  required_providers {
    wiz = {
      source = "iancrichardson/wiz"
      version = "0.4.0"
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

# Create a basic connector
resource "wiz_connector" "azure_basic" {
  name = "Azure Basic Connector"
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

# Create a connector with log monitor using OAuth
resource "wiz_connector" "azure_with_log_monitor" {
  name = "Azure Connector with Log Monitor"
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
    auditLogMonitorEnabled = true
    azureMonitorConfig = {
      eventHub = {
        connectionMethod = "OAUTH_SINGLE_BY_NAME"
        name = "wiz-cloud-events-hub"
        namespace = "wiz-cloud-events-namespace"
        namespaceTag = ""
      }
    }
    scheduledSecurityToolScanningSettings = {
      enabled = true
      publicBucketsScanningEnabled = false
    }
  })
}

# Output the basic connector ID
output "basic_connector_id" {
  value = wiz_connector.azure_basic.id
}

# Output the log monitor connector ID
output "log_monitor_connector_id" {
  value = wiz_connector.azure_with_log_monitor.id
}

# Output the configuration test result
output "config_test_success" {
  value = data.wiz_connector_config.test.success
}
