package resources

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jimmyflatting/terraform-provider-litellm/internal/provider"
)

func TestAccResourceKey_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { provider.TestAccPreCheck(t) },
		ProviderFactories: provider.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceKeyConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"litellm_key.test", "key_alias", "test-key"),
					resource.TestCheckResourceAttr(
						"litellm_key.test", "team_id", "test-team"),
					resource.TestCheckResourceAttr(
						"litellm_key.test", "max_budget", "100"),
				),
			},
			// Test update
			{
				Config: testAccResourceKeyConfig_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"litellm_key.test", "key_alias", "test-key-updated"),
					resource.TestCheckResourceAttr(
						"litellm_key.test", "team_id", "test-team-2"),
					resource.TestCheckResourceAttr(
						"litellm_key.test", "max_budget", "200"),
				),
			},
			// Import test
			{
				ResourceName:      "litellm_key.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceKey_fullConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { provider.TestAccPreCheck(t) },
		ProviderFactories: provider.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceKeyConfig_full(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"litellm_key.full", "key_alias", "full-test-key"),
					resource.TestCheckResourceAttr(
						"litellm_key.full", "team_id", "full-test-team"),
					resource.TestCheckResourceAttr(
						"litellm_key.full", "max_budget", "1000"),
					resource.TestCheckResourceAttr(
						"litellm_key.full", "models.#", "2"),
					resource.TestCheckResourceAttr(
						"litellm_key.full", "models.0", "gpt-4"),
					resource.TestCheckResourceAttr(
						"litellm_key.full", "models.1", "gpt-3.5-turbo"),
					resource.TestCheckResourceAttr(
						"litellm_key.full", "metadata.environment", "production"),
					resource.TestCheckResourceAttr(
						"litellm_key.full", "metadata.description", "Full test key"),
				),
			},
		},
	})
}

func testAccResourceKeyConfig_basic() string {
	return fmt.Sprintf(`
resource "litellm_key" "test" {
  key_alias  = "test-key"
  team_id    = "test-team"
  max_budget = 100
  models     = ["test-model"]
}
`)
}

func testAccResourceKeyConfig_update() string {
	return fmt.Sprintf(`
resource "litellm_key" "test" {
  key_alias  = "test-key-updated"
  team_id    = "test-team-2"
  max_budget = 200
  models     = ["test-model"]
}
`)
}

func testAccResourceKeyConfig_full() string {
	return fmt.Sprintf(`
resource "litellm_key" "full" {
  key_alias  = "full-test-key"
  team_id    = "full-test-team"
  max_budget = 1000
  models     = ["gpt-4", "gpt-3.5-turbo"]
  metadata = {
    environment = "production"
    description = "Full test key"
  }
}
`)
}
