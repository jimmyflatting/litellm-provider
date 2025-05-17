package validation

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// StringNotEmpty validates that a string value is not empty
func StringNotEmpty(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %s to be string", k)}
	}
	if v == "" {
		return nil, []error{fmt.Errorf("%s cannot be empty", k)}
	}
	return nil, nil
}

// StringMinLength validates that a string value has a minimum length
func StringMinLength(minLength int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) ([]string, []error) {
		v, ok := i.(string)
		if !ok {
			return nil, []error{fmt.Errorf("expected type of %s to be string", k)}
		}
		if len(v) < minLength {
			return nil, []error{fmt.Errorf("%s must be at least %d characters", k, minLength)}
		}
		return nil, nil
	}
}

// FloatGreaterThanOrEqual validates that a float64 value is greater than or equal to a minimum value
func FloatGreaterThanOrEqual(min float64) schema.SchemaValidateFunc {
	return func(i interface{}, k string) ([]string, []error) {
		v, ok := i.(float64)
		if !ok {
			return nil, []error{fmt.Errorf("expected type of %s to be float64", k)}
		}
		if v < min {
			return nil, []error{fmt.Errorf("%s cannot be less than %f", k, min)}
		}
		return nil, nil
	}
}

// OneOf validates that a string value is one of a set of values
func OneOf(valid ...string) schema.SchemaValidateFunc {
	return func(i interface{}, k string) ([]string, []error) {
		v, ok := i.(string)
		if !ok {
			return nil, []error{fmt.Errorf("expected type of %s to be string", k)}
		}
		for _, str := range valid {
			if v == str {
				return nil, nil
			}
		}
		return nil, []error{fmt.Errorf("expected %s to be one of %v, got %s", k, valid, v)}
	}
}
