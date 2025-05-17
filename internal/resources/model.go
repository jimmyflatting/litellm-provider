package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jimmyflatting/terraform-provider-litellm/internal/client"
	"github.com/jimmyflatting/terraform-provider-litellm/internal/validation"
)

func ResourceModel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceModelCreate,
		ReadContext:   resourceModelRead,
		UpdateContext: resourceModelUpdate,
		DeleteContext: resourceModelDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The name of the model",
				ValidateFunc: validation.StringMinLength(3),
			},
			"model_provider": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The provider of the model (e.g., 'openai', 'anthropic')",
				ValidateFunc: validation.OneOf(
					"openai",
					"anthropic",
					"azure",
					"cohere",
					"google",
					"replicate",
					"huggingface",
				),
			},
			"model_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the underlying model",
			},
			"api_base": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The base URL for API calls",
			},
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "API key for the model provider",
			},
			"metadata": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Additional metadata for the model",
			},
		},
	}
}

func resourceModelCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	model := &client.Model{
		Name:          d.Get("name").(string),
		ModelProvider: d.Get("model_provider").(string),
		ModelName:     d.Get("model_name").(string),
		APIBase:       d.Get("api_base").(string),
		APIKey:        d.Get("api_key").(string),
	}

	if v, ok := d.GetOk("metadata"); ok {
		metadata := make(map[string]string)
		for k, v := range v.(map[string]interface{}) {
			metadata[k] = v.(string)
		}
		model.Metadata = metadata
	}

	if err := c.CreateModel(model); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(model.Name)

	return resourceModelRead(ctx, d, m)
}

func resourceModelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	model, err := c.GetModel(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if model == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", model.Name)
	d.Set("model_provider", model.ModelProvider)
	d.Set("model_name", model.ModelName)
	d.Set("api_base", model.APIBase)
	d.Set("metadata", model.Metadata)
	// Don't set api_key as it's sensitive and not returned by the API

	return nil
}

func resourceModelUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	model := &client.Model{
		Name:          d.Id(),
		ModelProvider: d.Get("model_provider").(string),
		ModelName:     d.Get("model_name").(string),
		APIBase:       d.Get("api_base").(string),
		APIKey:        d.Get("api_key").(string),
	}

	if v, ok := d.GetOk("metadata"); ok {
		metadata := make(map[string]string)
		for k, v := range v.(map[string]interface{}) {
			metadata[k] = v.(string)
		}
		model.Metadata = metadata
	}

	if err := c.UpdateModel(model); err != nil {
		return diag.FromErr(err)
	}

	return resourceModelRead(ctx, d, m)
}

func resourceModelDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	if err := c.DeleteModel(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
