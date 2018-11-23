package atgen

import (
	"testing"
)

func TestGetVersionsOfTestFunc(t *testing.T) {
	parsed, err := parseYaml([]byte(yamlTestFuncPerAPIVersion))
	if err != nil {
		t.Fatal(err)
	}
	testFuncs := convertToTestFuncs(parsed)
	versions := getVersions(testFuncs[0])

	if versions[0] != "v1beta1" {
		t.Fatalf("versions[0] should be v1beta1")
	}

	if versions[1] != "v1beta2" {
		t.Fatalf("versions[1] should be v1beta2")
	}

	if versions[2] != "v1" {
		t.Fatalf("versions[2] should be v1")
	}
}

func TestGetVersionsOfTest(t *testing.T) {
	parsed, err := parseYaml([]byte(yamlTestPerAPIVersion))
	if err != nil {
		t.Fatal(err)
	}
	testFuncs := convertToTestFuncs(parsed)
	versions := getVersions(testFuncs[0])

	if len(versions) != 3 {
		t.Logf("%#v", versions)
		t.Fatalf("len(versions) should be 3, but %d", len(versions))
	}

	if versions[0] != "v1beta1" {
		t.Fatalf("versions[0] should be v1beta1")
	}

	if versions[1] != "v1beta2" {
		t.Fatalf("versions[1] should be v1beta2")
	}

	if versions[2] != "v1" {
		t.Fatalf("versions[2] should be v1")
	}
}

func TestGetVersionsOfTestFuncAndTest(t *testing.T) {
	parsed, err := parseYaml([]byte(yamlTestFuncAndTestPerAPIVersion))
	if err != nil {
		t.Fatal(err)
	}
	testFuncs := convertToTestFuncs(parsed)
	versions := getVersions(testFuncs[0])

	if len(versions) != 3 {
		t.Logf("%#v", versions)
		t.Fatalf("len(versions) should be 3, but %d", len(versions))
	}

	if versions[0] != "v1beta1" {
		t.Fatalf("versions[0] should be v1beta1")
	}

	if versions[1] != "v1beta2" {
		t.Fatalf("versions[1] should be v1beta2")
	}

	if versions[2] != "v1" {
		t.Fatalf("versions[2] should be v1")
	}
}

func TestFilterTestFunc(t *testing.T) {
	parsed, err := parseYaml([]byte(yamlTestFuncPerAPIVersion))
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
