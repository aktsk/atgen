package main

import (
	"testing"
)

func checkCompare(t *testing.T, actual, expected interface{}) {
	t.Helper()

	if actual == nil {
		t.Fatalf("Expected response should include %#v, but actually %#v", expected, actual)
	}

	a := actual.(map[string]interface{})
	for k, v := range expected.(map[string]interface{}) {
		switch p := v.(type) {
		case int:
			v = float64(p)
			if a[k] != v {
				t.Fatalf("Expected response parameter %v should be %#v, but actually %#v", k, v, a[k])
			}
		case string, bool:
			if a[k] != v {
				t.Fatalf("Expected response parameter %v should be %#v, but actually %#v", k, v, a[k])
			}
		case []map[string]interface{}:
			for i, e := range p {
				checkCompare(t, a[k].([]interface{})[i], e)
			}
		case map[string]interface{}:
			checkCompare(t, a[k], v)
		}
	}
}
