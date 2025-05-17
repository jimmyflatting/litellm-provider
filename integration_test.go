package main

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jimmyflatting/terraform-provider-litellm/internal/provider"
)

func TestIntegration_ModelAndKey(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { provider.TestAccPreCheck(t) },
		ProviderFactories: provider.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testIntegrationConfig_modelAndKey(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"litellm_model.integration", "name", "integration-model"),
					resource.TestCheckResourceAttr(
						"litellm_key.integration", "key_alias", "integration-key"),
					resource.TestCheckResourceAttr(
						"litellm_key.integration", "models.0", "integration-model"),
				),
			},
			{
				Config: testIntegrationConfig_updateModelAndKey(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"litellm_model.integration", "name", "integration-model-updated"),
					resource.TestCheckResourceAttr(
						"litellm_key.integration", "key_alias", "integration-key-updated"),
					resource.TestCheckResourceAttr(
						"litellm_key.integration", "models.0", "integration-model-updated"),
				),
			},
		},
	})
}

func testIntegrationConfig_modelAndKey() string {
	return `
resource "litellm_model" "integration" {
  name           = "integration-model"
  model_provider = "openai"
  model_name     = "gpt-4"
  api_base       = "https://api.openai.com/v1"
  metadata = {
    description = "Integration test model"
  }
}

resource "litellm_key" "integration" {
  key_alias  = "integration-key"
  team_id    = "integration-team"
  max_budget = 500
  models     = [litellm_model.integration.name]
  metadata = {
    description = "Integration test key"
  }
}
`
}

func testIntegrationConfig_updateModelAndKey() string {
	return `
resource "litellm_model" "integration" {
  name           = "integration-model-updated"
  model_provider = "openai"
  model_name     = "gpt-4"
  api_base       = "https://api.openai.com/v1"
  metadata = {
    description = "Updated integration test model"
  }
}

resource "litellm_key" "integration" {
  key_alias  = "integration-key-updated"
  team_id    = "integration-team"
  max_budget = 1000
  models     = [litellm_model.integration.name]
  metadata = {
    description = "Updated integration test key"
  }
}
`
}
