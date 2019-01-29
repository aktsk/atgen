package atgen

import (
	"testing"
)

func convertToTestFuncsHelper(t *testing.T, parsed []map[interface{}]interface{}) TestFuncs {
	t.Helper()

	funcs, err := convertToTestFuncs(parsed)
	if err != nil {
		t.Fatal(err)
	}
	return funcs
}

func TestParseTestFuncPerAPIVersion(t *testing.T) {
	parsed, err := parseYaml([]byte(yamlTestFuncPerAPIVersion))
	if err != nil {
		t.Fatal(err)
	}

	testFuncs := convertToTestFuncsHelper(t, parsed)
	testFunc := testFuncs[0]

	if testFunc.Name != "TestFoo" {
		t.Fatal("testFunc.Name should be TestFoo")
	}

	if testFunc.APIVersions[0] != "v1beta1" {
		t.Fatal("testFunc.APIVersions[0] should be v1beta1")
	}

	if testFunc.Vars["key"] != "val" {
		t.Fatal(`testFunc.Vars["key"] should be val`)
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

	testFuncs := convertToTestFuncsHelper(t, parsed)
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

	testFuncs := convertToTestFuncsHelper(t, parsed)
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

	testFuncs := convertToTestFuncsHelper(t, parsed)
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
	if test.Path != "/{apiVersion}/user/1/foo" {
		t.Fatal("test.Path should be /{apiVersion}/user/1/foo")
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
	testFuncs := convertToTestFuncsHelper(t, parsed)
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
	testFuncs := convertToTestFuncsHelper(t, parsed)
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
	testFuncs := convertToTestFuncsHelper(t, parsed)
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

func TestAggregateRouterFunc(t *testing.T) {
	testCases := []struct {
		Name                string
		Yaml                string
		TemplatePackagePath string
		ExpectedPackagePath string
		ExpectedName        string
	}{
		{
			Name: "Set router with function name only",
			Yaml: `
- name: TestSetRouterWithFunctionNameOnly
  routerFunc: getRouter
  apiVersions:
    - v1
  tests:
    - path: /{apiVersion}/user
      method: get
      res:
        status: 200
        params:
          foo: bar
`,
			TemplatePackagePath: "github.com/foo/bar/template",
			ExpectedName:        "getRouter",
			ExpectedPackagePath: "github.com/foo/bar/template",
		},
		{
			Name: "Set router with relative path",
			Yaml: `
- name: TestSetRouterWithRelativePath
  routerFunc: ../handlers.GetRouter
  apiVersions:
    - v1
  tests:
    - path: /{apiVersion}/user
      method: get
      res:
        status: 200
        params:
          foo: bar
`,
			TemplatePackagePath: "github.com/foo/bar/template",
			ExpectedName:        "GetRouter",
			ExpectedPackagePath: "github.com/foo/bar/handlers",
		},
		{
			Name: "Set router with absolute path",
			Yaml: `
- name: TestSetRouterWithAbsolute
  routerFunc: github.com/foo/bar/handlers.GetRouter
  apiVersions:
    - v1
  tests:
    - path: /{apiVersion}/user
      method: get
      res:
        status: 200
        params:
          foo: bar
`,
			TemplatePackagePath: "github.com/foo/bar/template",
			ExpectedName:        "GetRouter",
			ExpectedPackagePath: "github.com/foo/bar/handlers",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			parsed, err := parseYaml([]byte(testCase.Yaml))
			if err != nil {
				t.Fatal(err)
			}
			testFuncs := convertToTestFuncsHelper(t, parsed)
			routerFuncs, _, err := aggregateRouterFunc(testFuncs, testCase.TemplatePackagePath)
			if err != nil {
				t.Fatal(err)
			}

			routerFunc := routerFuncs[0]

			if routerFunc.Name != testCase.ExpectedName {
				t.Errorf("Expected name is %s, but actually %s", testCase.ExpectedName, routerFunc.Name)
			}

			if routerFunc.PackagePath != testCase.ExpectedPackagePath {
				t.Errorf("Expected package path is %s, but actually %s", testCase.ExpectedPackagePath, routerFunc.PackagePath)
			}
		})
	}
}
