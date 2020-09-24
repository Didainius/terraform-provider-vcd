// +build vapp vm user nsxt ALL functional

package vcd

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// testCachedFieldValue structure with attached functions is useful for testing specific field value
// across different `resource.TestStep` in Terraform acceptance tests. One particular use case is
// to check whether MAC address does not change when a `vcd_vapp_vm` resource's network stack is
// updated (between different TestSteps).
type testCachedFieldValue struct {
	fieldValue string
}

// cacheTestResourceFieldValue has the same signature as builtin Terraform Test functions, however
// it is attached to a struct which allows to store a field value and then check against this value
// with 'testCheckCachedResourceFieldValue'
func (c *testCachedFieldValue) cacheTestResourceFieldValue(resource, field string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("resource not found: %s", resource)
		}

		value, exists := rs.Primary.Attributes[field]
		if !exists {
			return fmt.Errorf("field %s in resource %s does not exist", field, resource)
		}
		// Store the value in cache
		c.fieldValue = value
		return nil
	}
}

// testCheckCachedResourceFieldValue has the default signature of Terraform acceptance test
// functions, but is able to verify if the value is equal to previously cached value using
// 'cacheTestResourceFieldValue'. This allows to check if a particular field value changed across
// multiple resource.TestSteps.
func (c *testCachedFieldValue) testCheckCachedResourceFieldValue(resource, field string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("resource not found: %s", resource)
		}

		value, exists := rs.Primary.Attributes[field]
		if !exists {
			return fmt.Errorf("field %s in resource %s does not exist", field, resource)
		}

		if value != c.fieldValue {
			return fmt.Errorf("got '%s - %s' field value %s, expected: %s",
				resource, field, value, c.fieldValue)
		}

		return nil
	}
}

// testCheckMatchOutput allows to match output field with regexp
func testCheckMatchOutput(name string, r *regexp.Regexp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		rs, ok := ms.Outputs[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if !r.MatchString(rs.Value.(string)) {
			return fmt.Errorf(
				"Output '%s': expected %#v, got %#v", name, rs.Value, rs)
		}

		return nil
	}
}

// testCheckOutputNonEmpty checks that output field is not empty
func testCheckOutputNonEmpty(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		rs, ok := ms.Outputs[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Value.(string) == "" {
			return fmt.Errorf("Output '%s': expected '', got %#v", name, rs)
		}

		return nil
	}
}

// resourceFieldsEqual checks if secondObject has all the fields and their values set as the
// firstObject except `[]excludeFields`. This is very useful to check if data sources have all
// the same values as resources
func resourceFieldsEqual(firstObject, secondObject string, excludeFields []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource1, ok := s.RootModule().Resources[firstObject]
		if !ok {
			return fmt.Errorf("unable to find %s", firstObject)
		}

		resource2, ok := s.RootModule().Resources[secondObject]
		if !ok {
			return fmt.Errorf("unable to find %s", secondObject)
		}

		for fieldName := range resource1.Primary.Attributes {
			// Do not validate the fields marked for exclusion
			if stringInSlice(fieldName, excludeFields) {
				continue
			}

			if vcdTestVerbose {
				fmt.Printf("field %s %s (value %s) and %s (value %s))\n", fieldName, firstObject,
					resource1.Primary.Attributes[fieldName], secondObject, resource2.Primary.Attributes[fieldName])
			}
			if !reflect.DeepEqual(resource1.Primary.Attributes[fieldName], resource2.Primary.Attributes[fieldName]) {
				return fmt.Errorf("field %s differs in resources %s (value %s) and %s (value %s)",
					fieldName, firstObject, resource1.Primary.Attributes[fieldName], secondObject, resource2.Primary.Attributes[fieldName])
			}
		}
		return nil
	}
}
