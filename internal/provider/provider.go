package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jimmyflatting/terraform-provider-litellm/internal/client"
	"github.com/jimmyflatting/terraform-provider-litellm/internal/datasources"
	"github.com/jimmyflatting/terraform-provider-litellm/internal/resources"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	schema.DescriptionKind = schema.StringMarkdown
}

func New() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("LITELLM_API_KEY", nil),
				Description: "API Key for authenticating with LiteLLM.",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://api.litellm.io",
				Description: "Base URL for the LiteLLM API.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"litellm_model": resources.ResourceModel(),
			"litellm_key":   resources.ResourceKey(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"litellm_model": datasources.DataSourceModel(),
			"litellm_key":   datasources.DataSourceKey(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

type Config struct {
	APIKey   string
	Endpoint string
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("api_key").(string)
	endpoint := d.Get("endpoint").(string)

	client := client.NewClient(apiKey, endpoint)

	return client, nil
}
