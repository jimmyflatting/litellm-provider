package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var TestAccProvider *schema.Provider
var TestAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	TestAccProvider = New()
	TestAccProviderFactories = map[string]func() (*schema.Provider, error){
		"litellm": func() (*schema.Provider, error) {
			return TestAccProvider, nil
		},
	}
}

func TestAccPreCheck(t *testing.T) {
	if v := os.Getenv("LITELLM_API_KEY"); v == "" {
		t.Fatal("LITELLM_API_KEY must be set for acceptance tests")
	}
}
