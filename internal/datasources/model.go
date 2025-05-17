package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jimmyflatting/terraform-provider-litellm/internal/client"
)

func DataSourceModel() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceModelRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the model",
			},
			"model_provider": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The provider of the model (e.g., 'openai', 'anthropic')",
			},
			"model_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the underlying model",
			},
			"api_base": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The base URL for API calls",
			},
			"metadata": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Additional metadata for the model",
			},
		},
	}
}

func dataSourceModelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	name := d.Get("name").(string)

	model, err := c.GetModel(name)
	if err != nil {
		return diag.FromErr(err)
	}

	if model == nil {
		return diag.Errorf("model %s not found", name)
	}

	d.SetId(model.Name)
	d.Set("model_provider", model.ModelProvider)
	d.Set("model_name", model.ModelName)
	d.Set("api_base", model.APIBase)
	d.Set("metadata", model.Metadata)

	return nil
}
