package atgen

import (
	"testing"
)

func TestFilterTestFunc(t *testing.T) {
	parsed, err := parseYaml([]byte(yamlTestFuncPerAPIVersion))
	if err != nil {
		t.Fatal(err)
	}
	tfuncs := filterTestFuncs(convertToTestFuncs(parsed))

	if tfuncs["v1beta1"][0].Vars["adminAPIKey"] != "test" {
		t.Fatalf(`tfuncs["v1beta1"][0].Vars["adminAPIKey"] should be test`)
	}

	v1beta1 := tfuncs["v1beta1"][0].Tests[0].(Test)
	if v1beta1.Path != "/v1beta1/money" {
		t.Fatalf("path should be /v1beta1/money, but %s", v1beta1.Path)
	}

	v1beta2 := tfuncs["v1beta2"][0].Tests[0].(Test)
	if v1beta2.Path != "/v1beta2/money" {
		t.Fatalf("path should be /v1beta2/money, but %s", v1beta2.Path)
	}
}

func TestFilterTest(t *testing.T) {
	parsed, err := parseYaml([]byte(yamlTestPerAPIVersion))
	if err != nil {
		t.Fatal(err)
	}
	tfuncs := filterTestFuncs(convertToTestFuncs(parsed))

	v1beta1 := tfuncs["v1beta1"][0].Tests[0].(Test)
	if v1beta1.Path != "/v1beta1/money" {
		t.Fatalf("path should be /v1beta1/money, but %s", v1beta1.Path)
	}

	v1beta2 := tfuncs["v1beta2"][0].Tests[0].(Test)
	if v1beta2.Path != "/v1beta2/money" {
		t.Fatalf("path should be /v1beta2/money, but %s", v1beta2.Path)
	}
}

func TestFilterTestFuncAndTest(t *testing.T) {
	parsed, err := parseYaml([]byte(yamlTestFuncAndTestPerAPIVersion))
	if err != nil {
		t.Fatal(err)
	}
	tfuncs := filterTestFuncs(convertToTestFuncs(parsed))

	v1beta1 := tfuncs["v1beta1"][0].Tests[0].(Test)
	if v1beta1.Path != "/v1beta1/money" {
		t.Fatalf("path should be /v1beta1/money, but %s", v1beta1.Path)
	}

	n := len(tfuncs["v1beta1"][0].Tests)
	if n != 1 {
		t.Fatalf("number of tests of v1beta1 should be 1, but %d", n)
	}

	v1beta2 := tfuncs["v1beta2"][0].Tests[0].(Test)
	if v1beta2.Path != "/v1beta2/money" {
		t.Fatalf("path should be /v1beta2/money, but %s", v1beta2.Path)
	}

	v1 := tfuncs["v1"][0].Tests[0].(Test)
	if v1.Path != "/v1/user" {
		t.Fatalf("path should be /v1/user, but %s", v1.Path)
	}
}
