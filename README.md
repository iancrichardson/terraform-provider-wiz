# Terraform Provider for Wiz

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
      version = "0.1.0"
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

### Installing the provider locally

To install the provider locally for testing:

```sh
mkdir -p ~/.terraform.d/plugins/iancrichardson/wiz/0.1.0/[OS]_[ARCH]/
cp terraform-provider-wiz ~/.terraform.d/plugins/iancrichardson/wiz/0.1.0/[OS]_[ARCH]/
```

Replace `[OS]_[ARCH]` with your system's OS and architecture (e.g., `darwin_amd64` for macOS on Intel, `darwin_arm64` for macOS on Apple Silicon).

## License

[MIT](LICENSE)
