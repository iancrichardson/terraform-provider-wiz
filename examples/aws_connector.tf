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

# Create a basic AWS connector
resource "wiz_connector" "aws_basic" {
  name = "AWS Basic Connector"
  type = "aws"
  
  auth_params = jsonencode({
    roleArn = "arn:aws:iam::123456789012:role/WizConnectorRole"
    externalId = "wiz-external-id"
  })
  
  extra_config = jsonencode({
    region = "us-east-1"
    scheduledSecurityToolScanningSettings = {
      enabled = true
      publicBucketsScanningEnabled = true
    }
  })
}

# Create an AWS connector with CloudTrail log monitoring
resource "wiz_connector" "aws_with_log_monitor" {
  name = "AWS Connector with Log Monitor"
  type = "aws"
  
  auth_params = jsonencode({
    roleArn = "arn:aws:iam::123456789012:role/WizConnectorRole"
    externalId = "wiz-external-id"
  })
  
  extra_config = jsonencode({
    region = "us-east-1"
    auditLogMonitorEnabled = true
    cloudtrailConfig = {
      s3 = {
        bucketName = "wiz-cloudtrail-logs"
        prefix = "AWSLogs/"
        roleArn = "arn:aws:iam::123456789012:role/WizCloudTrailRole"
      }
    }
    scheduledSecurityToolScanningSettings = {
      enabled = true
      publicBucketsScanningEnabled = true
    }
  })
}

# Output the AWS connector IDs
output "aws_basic_connector_id" {
  value = wiz_connector.aws_basic.id
}

output "aws_log_monitor_connector_id" {
  value = wiz_connector.aws_with_log_monitor.id
}
