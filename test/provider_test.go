package test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jimmyflatting/terraform-provider-litellm/internal/provider"
)

var testAccProvider *schema.Provider
var testAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	testAccProvider = provider.New()
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"litellm": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("LITELLM_API_KEY"); v == "" {
		t.Fatal("LITELLM_API_KEY must be set for acceptance tests")
	}
}

func TestAccResourceModel_InvalidConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceModelConfig_invalidProvider(),
				ExpectError: regexp.MustCompile(`Invalid model provider`),
			},
			{
				Config:      testAccResourceModelConfig_invalidModelName(),
				ExpectError: regexp.MustCompile(`Invalid model name`),
			},
			{
				Config:      testAccResourceModelConfig_invalidTimeout(),
				ExpectError: regexp.MustCompile(`Invalid timeout value`),
			},
		},
	})
}

func testAccResourceModelConfig_invalidProvider() string {
	return `
resource "litellm_model" "invalid" {
  name           = "invalid-model"
  model_provider = "invalid-provider"  # Invalid provider
  model_name     = "gpt-4"
  api_base       = "https://api.openai.com/v1"
}
`
}

func testAccResourceModelConfig_invalidModelName() string {
	return `
resource "litellm_model" "invalid" {
  name           = "invalid-model"
  model_provider = "openai"
  model_name     = ""  # Empty model name
  api_base       = "https://api.openai.com/v1"
}
`
}

func testAccResourceModelConfig_invalidTimeout() string {
	return `
resource "litellm_model" "invalid" {
  name           = "invalid-model"
  model_provider = "openai"
  model_name     = "gpt-4"
  api_base       = "https://api.openai.com/v1"
  timeout        = -1  # Negative timeout
}
`
}

func TestAccResourceKey_InvalidConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceKeyConfig_invalidAlias(),
				ExpectError: regexp.MustCompile(`Invalid key alias`),
			},
			{
				Config:      testAccResourceKeyConfig_invalidTeamId(),
				ExpectError: regexp.MustCompile(`Invalid team ID`),
			},
			{
				Config:      testAccResourceKeyConfig_invalidBudget(),
				ExpectError: regexp.MustCompile(`Invalid budget value`),
			},
			{
				Config:      testAccResourceKeyConfig_invalidModels(),
				ExpectError: regexp.MustCompile(`Invalid models configuration`),
			},
		},
	})
}

func testAccResourceKeyConfig_invalidAlias() string {
	return `
resource "litellm_key" "invalid" {
  key_alias  = ""  # Empty alias
  team_id    = "test-team"
  max_budget = 100
  models     = ["test-model"]
}
`
}

func testAccResourceKeyConfig_invalidTeamId() string {
	return `
resource "litellm_key" "invalid" {
  key_alias  = "test-key"
  team_id    = ""  # Empty team ID
  max_budget = 100
  models     = ["test-model"]
}
`
}

func testAccResourceKeyConfig_invalidBudget() string {
	return `
resource "litellm_key" "invalid" {
  key_alias  = "test-key"
  team_id    = "test-team"
  max_budget = -100  # Negative budget
  models     = ["test-model"]
}
`
}

func testAccResourceKeyConfig_invalidModels() string {
	return `
resource "litellm_key" "invalid" {
  key_alias  = "test-key"
  team_id    = "test-team"
  max_budget = 100
  models     = []  # Empty models list
}
`
}

func TestAccResourceIntegration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIntegrationConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"litellm_model.test", "name", "test-model"),
					resource.TestCheckResourceAttr(
						"litellm_key.test", "key_alias", "test-key"),
					resource.TestCheckResourceAttr(
						"litellm_key.test", "models.0", "test-model"),
				),
			},
		},
	})
}

func testAccResourceIntegrationConfig() string {
	return `
resource "litellm_model" "test" {
  name           = "test-model"
  model_provider = "openai"
  model_name     = "gpt-4"
  api_base       = "https://api.openai.com/v1"
}

resource "litellm_key" "test" {
  key_alias  = "test-key"
  team_id    = "test-team"
  max_budget = 100
  models     = [litellm_model.test.name]
}
`
}
