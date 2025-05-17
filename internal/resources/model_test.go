package resources

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jimmyflatting/terraform-provider-litellm/internal/provider"
)

func TestAccResourceModel_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { provider.TestAccPreCheck(t) },
		ProviderFactories: provider.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceModelConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"litellm_model.test", "name", "test-model"),
					resource.TestCheckResourceAttr(
						"litellm_model.test", "model_provider", "openai"),
					resource.TestCheckResourceAttr(
						"litellm_model.test", "model_name", "gpt-4"),
				),
			},
			// Test update
			{
				Config: testAccResourceModelConfig_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"litellm_model.test", "name", "test-model-updated"),
					resource.TestCheckResourceAttr(
						"litellm_model.test", "model_provider", "openai"),
					resource.TestCheckResourceAttr(
						"litellm_model.test", "model_name", "gpt-3.5-turbo"),
					resource.TestCheckResourceAttr(
						"litellm_model.test", "metadata.description", "Updated test model"),
				),
			},
			// Import test
			{
				ResourceName:      "litellm_model.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceModel_fullConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { provider.TestAccPreCheck(t) },
		ProviderFactories: provider.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceModelConfig_full(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"litellm_model.full", "name", "full-test-model"),
					resource.TestCheckResourceAttr(
						"litellm_model.full", "model_provider", "openai"),
					resource.TestCheckResourceAttr(
						"litellm_model.full", "model_name", "gpt-4"),
					resource.TestCheckResourceAttr(
						"litellm_model.full", "timeout", "30"),
					resource.TestCheckResourceAttr(
						"litellm_model.full", "metadata.description", "Full test model"),
					resource.TestCheckResourceAttr(
						"litellm_model.full", "metadata.environment", "testing"),
					resource.TestCheckResourceAttr(
						"litellm_model.full", "metadata.team", "qa"),
				),
			},
		},
	})
}

func testAccResourceModelConfig_basic() string {
	return fmt.Sprintf(`
resource "litellm_model" "test" {
  name           = "test-model"
  model_provider = "openai"
  model_name     = "gpt-4"
  api_base       = "https://api.openai.com/v1"
  metadata = {
    description = "Test model"
    team        = "testing"
  }
}
`)
}

func testAccResourceModelConfig_update() string {
	return fmt.Sprintf(`
resource "litellm_model" "test" {
  name           = "test-model-updated"
  model_provider = "openai"
  model_name     = "gpt-3.5-turbo"
  api_base       = "https://api.openai.com/v1"
  metadata = {
    description = "Updated test model"
    team        = "testing"
  }
}
`)
}

func testAccResourceModelConfig_full() string {
	return fmt.Sprintf(`
resource "litellm_model" "full" {
  name           = "full-test-model"
  model_provider = "openai"
  model_name     = "gpt-4"
  api_base       = "https://api.openai.com/v1"
  timeout        = 30
  metadata = {
    description  = "Full test model"
    environment  = "testing"
    team         = "qa"
  }
}
`)
}
