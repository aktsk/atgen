package atgen

import (
	"testing"
)

func TestParse(t *testing.T) {
	parsed, err := parseYaml([]byte(yamlString))
	if err != nil {
		t.Fatal(err)
	}

	testFuncs := convertToTestFuncs(parsed)
	testFunc := testFuncs[0]

	if testFunc.Name != "TestFoo" {
		t.Fatal("testFunc.Name should be TestFoo")
	}

	if testFunc.ApiVersion[0] != "v1beta1" {
		t.Fatal("testFunc.ApiVersion[0] should be v1beta1")
	}

	test := testFunc.Tests[0].(Test)
	if test.ApiVersion[0] != "v1beta1" {
		t.Fatal("test.ApiVersion[0]")
	}
}

var yamlString = `
- name: TestFoo
  apiVersion:
    - v1beta1
    - v1beta2
    - v1
  tests:
    - apiVersion:
        - v1beta1
        - v1beta2
        - v1
      path: /{apiVersion}/money
      method: post
      req:
        params:
          moneyId: "1"
          priority: free
          currency: JPY
        headers:
          x-admin-api-key: test
      res:
        status: 201
        headers:
          location: ""
        prams:
          foo: bar
`
