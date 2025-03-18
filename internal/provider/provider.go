package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/iancrichardson/terraform-provider-wiz/internal/client"
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The client ID for API operations",
				DefaultFunc: schema.EnvDefaultFunc("WIZ_CLIENT_ID", nil),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The client secret for API operations",
				DefaultFunc: schema.EnvDefaultFunc("WIZ_CLIENT_SECRET", nil),
			},
			"api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL of the Wiz GraphQL API",
				DefaultFunc: schema.EnvDefaultFunc("WIZ_API_URL", "https://api.eu1.demo.wiz.io/graphql"),
			},
			"auth_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL of the Wiz authentication endpoint",
				DefaultFunc: schema.EnvDefaultFunc("WIZ_AUTH_URL", "https://auth.demo.wiz.io/oauth/token"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"wiz_connector": resourceConnector(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"wiz_connector_config": dataSourceConnectorConfig(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	config := &client.Config{
		ClientID:     d.Get("client_id").(string),
		ClientSecret: d.Get("client_secret").(string),
		APIURL:       d.Get("api_url").(string),
		AuthURL:      d.Get("auth_url").(string),
	}

	c, err := client.NewClient(config)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return c, diags
}
