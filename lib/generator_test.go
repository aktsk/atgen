package atgen

import (
	"testing"
)

func parseAndConvertToTestFuncs(t *testing.T, yaml string) map[string]TestFuncs {
	t.Helper()

	parsed, err := parseYaml([]byte(yaml))
	if err != nil {
		t.Fatal(err)
	}

	testFuncs, err := convertToTestFuncs(parsed)
	if err != nil {
		t.Fatal(err)
	}

	return filterTestFuncs(testFuncs)
}

func TestFilterTestFunc(t *testing.T) {
	tfuncs := parseAndConvertToTestFuncs(t, yamlTestFuncPerAPIVersion)

	if tfuncs["v1beta1"][0].Vars["key"] != "val" {
		t.Fatalf(`tfuncs["v1beta1"][0].Vars["key"] should be val`)
	}

	v1beta1 := tfuncs["v1beta1"][0].Tests[0].(Test)
	if v1beta1.Path != "/v1beta1/user" {
		t.Fatalf("path should be /v1beta1/user, but %s", v1beta1.Path)
	}

	v1beta2 := tfuncs["v1beta2"][0].Tests[0].(Test)
	if v1beta2.Path != "/v1beta2/user" {
		t.Fatalf("path should be /v1beta2/user, but %s", v1beta2.Path)
	}
}

func TestFilterTest(t *testing.T) {
	tfuncs := parseAndConvertToTestFuncs(t, yamlTestPerAPIVersion)

	v1beta1 := tfuncs["v1beta1"][0].Tests[0].(Test)
	if v1beta1.Path != "/v1beta1/user" {
		t.Fatalf("path should be /v1beta1/user, but %s", v1beta1.Path)
	}

	v1beta2 := tfuncs["v1beta2"][0].Tests[0].(Test)
	if v1beta2.Path != "/v1beta2/user" {
		t.Fatalf("path should be /v1beta2/user, but %s", v1beta2.Path)
	}
}

func TestFilterTestFuncAndTest(t *testing.T) {
	tfuncs := parseAndConvertToTestFuncs(t, yamlTestFuncAndTestPerAPIVersion)

	v1beta1 := tfuncs["v1beta1"][0].Tests[0].(Test)
	if v1beta1.Path != "/v1beta1/user" {
		t.Fatalf("path should be /v1beta1/user, but %s", v1beta1.Path)
	}

	n := len(tfuncs["v1beta1"][0].Tests)
	if n != 1 {
		t.Fatalf("number of tests of v1beta1 should be 1, but %d", n)
	}

	v1beta2 := tfuncs["v1beta2"][0].Tests[0].(Test)
	if v1beta2.Path != "/v1beta2/user" {
		t.Fatalf("path should be /v1beta2/user, but %s", v1beta2.Path)
	}

	v1 := tfuncs["v1"][0].Tests[0].(Test)
	if v1.Path != "/v1/user" {
		t.Fatalf("path should be /v1/user, but %s", v1.Path)
	}
}

func TestFilterTestFuncAndSubtests(t *testing.T) {
	tfuncs := parseAndConvertToTestFuncs(t, yamlTestFuncWithSubtests)

	subtest := tfuncs["v1"][0].Tests[2].(Subtests)[0]
	if subtest.Name != "SubFoo" {
		t.Fatal("subtest.Name should be SubFoo")
	}

	if subtest.Tests[0].Path != "/v1/user/1/foo" {
		t.Fatal("subtest.Tests[0].Path should be /v1/user/1/foo")
	}
}
