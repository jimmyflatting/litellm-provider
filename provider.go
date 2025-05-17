package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
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
			// Resources will be added here
		},
		DataSourcesMap: map[string]*schema.Resource{
			// Data sources will be added here
		},
	}
}
