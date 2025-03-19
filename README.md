# Terraform Provider for Wiz

[![CI](https://github.com/iancrichardson/terraform-provider-wiz/actions/workflows/ci.yml/badge.svg)](https://github.com/iancrichardson/terraform-provider-wiz/actions/workflows/ci.yml)
[![Release](https://github.com/iancrichardson/terraform-provider-wiz/actions/workflows/release.yml/badge.svg)](https://github.com/iancrichardson/terraform-provider-wiz/actions/workflows/release.yml)

This Terraform provider allows you to manage Wiz connectors through Terraform. It provides resources and data sources for interacting with the Wiz API.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Go](https://golang.org/doc/install) >= 1.18 (to build the provider plugin)

## Building The Provider

1. Clone the repository
```sh
git clone https://github.com/iancrichardson/terraform-provider-wiz.git
```

2. Enter the provider directory
```sh
cd terraform-provider-wiz
```

3. Build the provider
```sh
go build -o terraform-provider-wiz
```

## Using the provider

To use the provider, add the following to your Terraform configuration:

```hcl
terraform {
  required_providers {
    wiz = {
      source = "iancrichardson/wiz"
      version = "0.4.0"
    }
  }
}

provider "wiz" {
  client_id     = "YOUR_CLIENT_ID"
  client_secret = "YOUR_CLIENT_SECRET"
  # Optional: Uncomment to override default endpoints
  # api_url       = "https://api.eu1.demo.wiz.io/graphql"
  # auth_url      = "https://auth.demo.wiz.io/oauth/token"
}
```

### Authentication

The provider supports the following authentication methods:

1. Static credentials
   ```hcl
   provider "wiz" {
     client_id     = "YOUR_CLIENT_ID"
     client_secret = "YOUR_CLIENT_SECRET"
   }
   ```

2. Environment variables
   ```sh
   export WIZ_CLIENT_ID="YOUR_CLIENT_ID"
   export WIZ_CLIENT_SECRET="YOUR_CLIENT_SECRET"
   ```

   ```hcl
   provider "wiz" {}
   ```

## Resources

### wiz_connector

The `wiz_connector` resource allows you to create and manage Wiz connectors.

#### Basic Connector Example

```hcl
resource "wiz_connector" "azure" {
  name = "Azure Connector"
  type = "azure"
  
  auth_params = jsonencode({
    isManagedIdentity = true
    subscriptionId    = "your-subscription-id"
    tenantId          = "your-tenant-id"
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
```

#### Connector with Log Monitoring

You can enable log monitoring for your connectors to collect audit logs from your cloud environments. The configuration varies by cloud provider:

##### Azure with OAuth Log Monitoring

```hcl
resource "wiz_connector" "azure_with_log_monitor" {
  name = "Azure Connector with Log Monitor"
  type = "azure"
  
  auth_params = jsonencode({
    isManagedIdentity = true
    subscriptionId    = "your-subscription-id"
    tenantId          = "your-tenant-id"
    environment       = "AzurePublicCloud"
  })
  
  extra_config = jsonencode({
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
```

##### AWS with CloudTrail Log Monitoring

```hcl
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
  })
}
```

##### GCP with Pub/Sub Log Monitoring

```hcl
resource "wiz_connector" "gcp_with_log_monitor" {
  name = "GCP Connector with Log Monitor"
  type = "gcp"
  
  auth_params = jsonencode({
    projectId = "my-gcp-project-id"
    serviceAccountKey = "REDACTED_SERVICE_ACCOUNT_KEY"
  })
  
  extra_config = jsonencode({
    auditLogMonitorEnabled = true
    auditLogsConfig = {
      pub_sub = {
        topicName = "projects/my-gcp-project-id/topics/wiz-cloud-events"
        subscriptionID = "wiz-cloud-events-sub"
      }
    }
  })
}
```

For more detailed examples, see the [examples directory](examples/).

## Data Sources

### wiz_connector_config

The `wiz_connector_config` data source allows you to test a connector configuration before creating it.

```hcl
data "wiz_connector_config" "test" {
  type = "azure"
  auth_params = jsonencode({
    isManagedIdentity = true
    subscriptionId    = "your-subscription-id"
    tenantId          = "your-tenant-id"
    environment       = "AzurePublicCloud"
  })
}

output "config_test_success" {
  value = data.wiz_connector_config.test.success
}
```

## Development

### Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Go](https://golang.org/doc/install) >= 1.18

### Building

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `build` command:
```sh
go build -o terraform-provider-wiz
```

### Installing the provider

#### Automated Installation (Terraform 0.13+)

The provider is available in the [Terraform Registry](https://registry.terraform.io/providers/iancrichardson/wiz/latest). To use it, simply specify it in your Terraform configuration:

```hcl
terraform {
  required_providers {
    wiz = {
      source  = "iancrichardson/wiz"
      version = "0.4.0"
    }
  }
}
```

Terraform will automatically download the provider when you run `terraform init`.

#### Manual Installation

If you prefer to install the provider manually:

1. Download the appropriate version for your platform from the [releases page](https://github.com/iancrichardson/terraform-provider-wiz/releases).
2. Extract the zip file.
3. Move the binary to the Terraform plugin directory:

```sh
mkdir -p ~/.terraform.d/plugins/github.com/iancrichardson/wiz/0.4.0/[OS]_[ARCH]/
mv terraform-provider-wiz_v0.4.0 ~/.terraform.d/plugins/github.com/iancrichardson/wiz/0.4.0/[OS]_[ARCH]/terraform-provider-wiz_v0.4.0
```

Replace `[OS]_[ARCH]` with your system's OS and architecture (e.g., `darwin_amd64` for macOS on Intel, `darwin_arm64` for macOS on Apple Silicon).

> **Note:** Starting with Terraform 0.13, the provider path includes `github.com/` in the path.

### Building from Source

To build the provider from source:

```sh
git clone https://github.com/iancrichardson/terraform-provider-wiz.git
cd terraform-provider-wiz
go build -o terraform-provider-wiz
```

Then install it locally:

```sh
mkdir -p ~/.terraform.d/plugins/github.com/iancrichardson/wiz/0.4.0/[OS]_[ARCH]/
cp terraform-provider-wiz ~/.terraform.d/plugins/github.com/iancrichardson/wiz/0.4.0/[OS]_[ARCH]/
```

## Version History

### 0.4.0

This version includes important improvements to handle sensitive fields during connector updates:

- **Fixed Update Issues**: Resolved issues when updating connectors with sensitive fields in `auth_params`
- **Improved Log Monitoring Support**: Enhanced support for enabling log monitoring with OAuth and other authentication methods
- **Better Error Handling**: Added more robust error handling for connector operations
- **Empty Auth Params**: When updating connectors, sensitive fields in `auth_params` are now excluded from the update operation to prevent API errors

This version is particularly useful if you're experiencing errors like:
```
Error: error updating connector: error updating connector: graphql: oops! an internal error has occurred
```

When updating connectors with log monitoring enabled.

## Contributing

Contributions are welcome! Here's how you can contribute to the project:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run the tests (`go test ./...`)
5. Commit your changes (`git commit -m 'Add some amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Development Workflow

1. Make changes to the provider code
2. Build the provider using `go build`
3. Test your changes locally by installing the provider in your local Terraform plugin directory
4. Run acceptance tests if applicable

### Releasing

Releases are automatically created when a new tag is pushed to the repository. To create a new release:

1. Update the version in `main.go`
2. Update the version in example files and README
3. Commit your changes
4. Create and push a new tag:
   ```sh
   git tag v0.x.0
   git push origin v0.x.0
   ```

The GitHub Actions workflow will automatically build the provider for all supported platforms and create a new release.

## License

[MIT](LICENSE)
