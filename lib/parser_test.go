package atgen

import (
	"testing"
)

func TestParseTestFuncPerAPIVersion(t *testing.T) {
	parsed, err := parseYaml([]byte(yamlTestFuncPerAPIVersion))
	if err != nil {
		t.Fatal(err)
	}

	testFuncs := convertToTestFuncs(parsed)
	testFunc := testFuncs[0]

	if testFunc.Name != "TestFoo" {
		t.Fatal("testFunc.Name should be TestFoo")
	}

	if testFunc.APIVersions[0] != "v1beta1" {
		t.Fatal("testFunc.APIVersions[0] should be v1beta1")
	}

	if testFunc.Vars["adminAPIKey"] != "test" {
		t.Fatal(`testFunc.Vars["adminAPIKey"] should be test`)
	}

	if testFunc.Vars["foo"].(map[string]interface{})["bar"] != "baz" {
		t.Fatal(`testFunc.Vars["foo"]["bar"] should be baz`)
	}

	test := testFunc.Tests[0].(Test)
	if test.Vars["foo"] != "bar" {
		t.Fatal(`test.Vars["foo"] should be bar`)
	}
}

func TestParseTestPerAPIVersion(t *testing.T) {
	parsed, err := parseYaml([]byte(yamlTestPerAPIVersion))
	if err != nil {
		t.Fatal(err)
	}

	testFuncs := convertToTestFuncs(parsed)
	testFunc := testFuncs[0]

	if testFunc.Name != "TestFoo" {
		t.Fatal("testFunc.Name should be TestFoo")
	}

	test := testFunc.Tests[0].(Test)
	if test.APIVersions[0] != "v1beta1" {
		t.Fatal("test.APIVersions[0]")
	}

	if test.Res.Params["foo"] != "bar" {
		t.Log(test.Res.Params["foo"])
		t.Fatal(`test.Res.Params["foo"] should be bar`)
	}
}

func TestParseTestFuncAndTestPerAPIVersion(t *testing.T) {
	parsed, err := parseYaml([]byte(yamlTestFuncAndTestPerAPIVersion))
	if err != nil {
		t.Fatal(err)
	}

	testFuncs := convertToTestFuncs(parsed)
	testFunc := testFuncs[0]

	if testFunc.Name != "TestFoo" {
		t.Fatal("testFunc.Name should be TestFoo")
	}

	if testFunc.APIVersions[0] != "v1beta1" {
		t.Fatal("testFunc.APIVersions[0] should be v1beta1")
	}

	test := testFunc.Tests[0].(Test)
	if test.APIVersions[0] != "v1beta1" {
		t.Fatal("test.APIVersions[0]")
	}
}

func TestParseTestFuncWithSubtests(t *testing.T) {
	parsed, err := parseYaml([]byte(yamlTestFuncWithSubtests))
	if err != nil {
		t.Fatal(err)
	}

	testFuncs := convertToTestFuncs(parsed)
	testFunc := testFuncs[0]

	if testFunc.Name != "TestWithSubtests" {
		t.Fatal("testFunc.Name should be TestFoo")
	}

	if testFunc.APIVersions[0] != "v1" {
		t.Fatal("testFunc.APIVersions[0] should be v1beta1")
	}

	subtest := testFunc.Tests[2].(Subtests)[0]
	if subtest.Name != "SubFoo" {
		t.Fatal("subtest.Name should be SubFoo")
	}

	test := subtest.Tests[0]
	if test.Path != "/{apiVersion}/money/1/sub" {
		t.Fatal("test.Path should be /{apiVersion}/money/1/sub")
	}

	if test.Method != "delete" {
		t.Fatal("test.Method should be delete")
	}

	if test.Vars["foo"] != "bar" {
		t.Fatal(`test.Vars["foo"] should be bar`)
	}
}

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
