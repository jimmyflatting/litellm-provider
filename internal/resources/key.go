package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jimmyflatting/terraform-provider-litellm/internal/client"
)

func ResourceKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeyCreate,
		ReadContext:   resourceKeyRead,
		UpdateContext: resourceKeyUpdate,
		DeleteContext: resourceKeyDelete,

		Schema: map[string]*schema.Schema{
			"key_alias": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Alias for the key",
				ValidateFunc: func(i interface{}, k string) ([]string, []error) {
					v, ok := i.(string)
					if !ok {
						return nil, []error{fmt.Errorf("expected type of %s to be string", k)}
					}
					if v == "" {
						return nil, []error{fmt.Errorf("%s cannot be empty", k)}
					}
					if len(v) < 3 {
						return nil, []error{fmt.Errorf("%s must be at least 3 characters", k)}
					}
					return nil, nil
				},
			},
			"team_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Team ID associated with the key",
				ValidateFunc: func(i interface{}, k string) ([]string, []error) {
					v, ok := i.(string)
					if !ok {
						return nil, []error{fmt.Errorf("expected type of %s to be string", k)}
					}
					if v == "" {
						return nil, []error{fmt.Errorf("%s cannot be empty", k)}
					}
					return nil, nil
				},
			},
			"models": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of models this key has access to",
			},
			"max_budget": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "Maximum budget allowed for this key",
				ValidateFunc: func(i interface{}, k string) ([]string, []error) {
					v, ok := i.(float64)
					if !ok {
						return nil, []error{fmt.Errorf("expected type of %s to be float64", k)}
					}
					if v < 0 {
						return nil, []error{fmt.Errorf("%s cannot be negative", k)}
					}
					return nil, nil
				},
			},
			"expires_at": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Expiration timestamp for the key",
			},
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The generated API key",
			},
		},
	}
}

func resourceKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	key := &client.Key{
		KeyAlias:  d.Get("key_alias").(string),
		TeamID:    d.Get("team_id").(string),
		MaxBudget: d.Get("max_budget").(float64),
		ExpiresAt: d.Get("expires_at").(string),
	}

	if v, ok := d.GetOk("models"); ok {
		models := make([]string, 0)
		for _, model := range v.([]interface{}) {
			models = append(models, model.(string))
		}
		key.Models = models
	}

	if err := c.CreateKey(key); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(key.KeyAlias)
	d.Set("key", key.Key)

	return resourceKeyRead(ctx, d, m)
}

func resourceKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	key, err := c.GetKey(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if key == nil {
		d.SetId("")
		return nil
	}

	d.Set("key_alias", key.KeyAlias)
	d.Set("team_id", key.TeamID)
	d.Set("models", key.Models)
	d.Set("max_budget", key.MaxBudget)
	d.Set("expires_at", key.ExpiresAt)
	// Note: The actual key value is only available during creation

	return nil
}

func resourceKeyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	key := &client.Key{
		KeyAlias:  d.Id(),
		TeamID:    d.Get("team_id").(string),
		MaxBudget: d.Get("max_budget").(float64),
		ExpiresAt: d.Get("expires_at").(string),
	}

	if v, ok := d.GetOk("models"); ok {
		models := make([]string, 0)
		for _, model := range v.([]interface{}) {
			models = append(models, model.(string))
		}
		key.Models = models
	}

	if err := c.UpdateKey(key); err != nil {
		return diag.FromErr(err)
	}

	return resourceKeyRead(ctx, d, m)
}

func resourceKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	if err := c.DeleteKey(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
