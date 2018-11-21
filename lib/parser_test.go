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

	if testFunc.APIVersion[0] != "v1beta1" {
		t.Fatal("testFunc.APIVersion[0] should be v1beta1")
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
	if test.APIVersion[0] != "v1beta1" {
		t.Fatal("test.APIVersion[0]")
	}
}

var yamlTestFuncPerAPIVersion = `
- name: TestFoo
  apiVersion:
    - v1beta1
    - v1beta2
    - v1
  tests:
    - path: /{apiVersion}/money
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

var yamlTestPerAPIVersion = `
- name: TestFoo
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
    - apiVersion:
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
