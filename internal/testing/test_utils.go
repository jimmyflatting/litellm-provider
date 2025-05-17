package testing

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestCheckResourceAttrWithRegex(name, key, pattern string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}

		value, ok := rs.Primary.Attributes[key]
		if !ok {
			return fmt.Errorf("attribute '%s' not found", key)
		}

		regex := regexp.MustCompile(pattern)
		if !regex.MatchString(value) {
			return fmt.Errorf(
				"attribute '%s' value '%s' does not match pattern '%s'",
				key, value, pattern,
			)
		}

		return nil
	}
}
