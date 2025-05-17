package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jimmyflatting/terraform-provider-litellm/internal/client"
)

func DataSourceKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKeyRead,

		Schema: map[string]*schema.Schema{
			"key_alias": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Alias for the key",
			},
			"team_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Team ID associated with the key",
			},
			"models": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of models this key has access to",
			},
			"max_budget": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Maximum budget allowed for this key",
			},
			"expires_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiration timestamp for the key",
			},
		},
	}
}

func dataSourceKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	keyAlias := d.Get("key_alias").(string)

	key, err := c.GetKey(keyAlias)
	if err != nil {
		return diag.FromErr(err)
	}

	if key == nil {
		return diag.Errorf("key with alias %s not found", keyAlias)
	}

	d.SetId(key.KeyAlias)
	d.Set("team_id", key.TeamID)
	d.Set("models", key.Models)
	d.Set("max_budget", key.MaxBudget)
	d.Set("expires_at", key.ExpiresAt)

	return nil
}
