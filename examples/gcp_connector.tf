terraform {
  required_providers {
    wiz = {
      source  = "iancrichardson/wiz"
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

# Create a basic GCP connector
resource "wiz_connector" "gcp_basic" {
  name = "GCP Basic Connector"
  type = "gcp"
  
  auth_params = jsonencode({
    projectId = "my-gcp-project-id"
    serviceAccountKey = <<EOT
{
  "type": "service_account",
  "project_id": "my-gcp-project-id",
  "private_key_id": "private-key-id",
  "private_key": "-----BEGIN PRIVATE KEY-----\nREDACTED\n-----END PRIVATE KEY-----\n",
  "client_email": "wiz-connector@my-gcp-project-id.iam.gserviceaccount.com",
  "client_id": "client-id",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/wiz-connector%40my-gcp-project-id.iam.gserviceaccount.com"
}
EOT
  })
  
  extra_config = jsonencode({
    projects = []
    excludedProjects = []
    includedFolders = []
    excludedFolders = []
    auditLogMonitorEnabled = false
    scheduledSecurityToolScanningSettings = {
      enabled = true
      publicBucketsScanningEnabled = true
    }
  })
}

# Create a GCP connector with Pub/Sub log monitoring
resource "wiz_connector" "gcp_with_log_monitor" {
  name = "GCP Connector with Log Monitor"
  type = "gcp"
  
  auth_params = jsonencode({
    projectId = "my-gcp-project-id"
    serviceAccountKey = <<EOT
{
  "type": "service_account",
  "project_id": "my-gcp-project-id",
  "private_key_id": "private-key-id",
  "private_key": "-----BEGIN PRIVATE KEY-----\nREDACTED\n-----END PRIVATE KEY-----\n",
  "client_email": "wiz-connector@my-gcp-project-id.iam.gserviceaccount.com",
  "client_id": "client-id",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/wiz-connector%40my-gcp-project-id.iam.gserviceaccount.com"
}
EOT
  })
  
  extra_config = jsonencode({
    projects = []
    excludedProjects = []
    includedFolders = []
    excludedFolders = []
    auditLogMonitorEnabled = true
    auditLogsConfig = {
      pub_sub = {
        topicName = "projects/my-gcp-project-id/topics/wiz-cloud-events"
        subscriptionID = "wiz-cloud-events-sub"
      }
    }
    scheduledSecurityToolScanningSettings = {
      enabled = true
      publicBucketsScanningEnabled = true
    }
  })
}

# Output the GCP connector IDs
output "gcp_basic_connector_id" {
  value = wiz_connector.gcp_basic.id
}

output "gcp_log_monitor_connector_id" {
  value = wiz_connector.gcp_with_log_monitor.id
}
