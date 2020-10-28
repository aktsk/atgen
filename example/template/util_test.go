package main

import (
	"fmt"
	"testing"
)

func compare(t *testing.T, actual, expected interface{}) error {
	if actual == nil {
		return fmt.Errorf("Expected response should include %#v, but actually %#v", expected, actual)
	}

	a := actual.(map[string]interface{})
	for k, v := range expected.(map[string]interface{}) {
		switch p := v.(type) {
		case int:
			v = float64(p)
			if a[k] != v {
				return fmt.Errorf("Expected response parameter %v should be %#v, but actually %#v", k, v, a[k])
			}
		case string, bool:
			if a[k] != v {
				return fmt.Errorf("Expected response parameter %v should be %#v, but actually %#v", k, v, a[k])
			}
		case []map[string]interface{}:
			for i, e := range p {
				err := compare(t, a[k].([]interface{})[i], e)
				if err != nil {
					return err
				}
			}
		case map[string]interface{}:
			err := compare(t, a[k], v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
