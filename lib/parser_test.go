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

var yamlTestFuncPerAPIVersion = `
- name: TestFoo
  apiVersions:
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
    - apiVersions:
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
    - apiVersions:
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

var yamlTestFuncAndTestPerAPIVersion = `
- name: TestFoo
  apiVersions:
    - v1beta1
    - v1beta2
    - v1
  tests:
    - apiVersions:
        - v1beta1
        - v1beta2
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
    - apiVersions:
        - v1
      path: /{apiVersion}/user
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
