package provider

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProvider *schema.Provider
var testAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	testAccProvider = New()
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"litellm": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := New().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("LITELLM_API_KEY"); v == "" {
		t.Fatal("LITELLM_API_KEY must be set for acceptance tests")
	}
}

func testAccProvider_Configure(t *testing.T) {
	raw := map[string]interface{}{
		"api_key": os.Getenv("LITELLM_API_KEY"),
	}

	diags := testAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("provider configuration failed: %v", diags)
	}
}
