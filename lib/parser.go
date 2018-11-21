package atgen

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

func (g *Generator) ParseYaml() error {
	buf, err := ioutil.ReadFile(g.Yaml)
	if err != nil {
		return err
	}

	parsed, err := parseYaml(buf)
	if err != nil {
		return err
	}

	g.TestFuncs = convertToTestFuncs(parsed)

	return nil
}

func parseYaml(buf []byte) ([]map[interface{}]interface{}, error) {
	var parsed []map[interface{}]interface{}
	err := yaml.Unmarshal(buf, &parsed)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

func convertToTestFuncs(parsed []map[interface{}]interface{}) TestFuncs {
	var testFuncs TestFuncs
	for _, p := range parsed {
		testFunc := TestFunc{}
		testFunc.Name = p["name"].(string)

		if p["apiVersion"] != nil {
			for _, v := range p["apiVersion"].([]interface{}) {
				testFunc.APIVersion = append(testFunc.APIVersion, v.(string))
			}
		}

		for _, t := range p["tests"].([]interface{}) {
			t := t.(map[interface{}]interface{})
			if t["path"] != nil {
				testFunc.Tests = append(testFunc.Tests, convertToTest(t))
			} else {
				testFunc.Tests = append(testFunc.Tests, convertToSubTests(t))
			}
		}
		testFuncs = append(testFuncs, testFunc)
	}
	return testFuncs
}

func convertToTest(t map[interface{}]interface{}) Test {
	var apiVersion []string
	if t["apiVersion"] != nil {
		for _, v := range t["apiVersion"].([]interface{}) {
			apiVersion = append(apiVersion, v.(string))
		}
	}

	return Test{
		APIVersion: apiVersion,
		Path:       t["path"].(string),
		Method:     t["method"].(string),
		Req:        convertToReq(t["req"]),
		Res:        convertToRes(t["res"]),
	}
}

func convertToSubTests(s map[interface{}]interface{}) SubTests {
	subTests := SubTests{}

	for _, t := range s["subtests"].([]interface{}) {
		test := convertToTest(t.(map[interface{}]interface{}))
		subTests = append(subTests, test)
	}

	return subTests
}

func convertToReq(r interface{}) Req {
	req := r.(map[interface{}]interface{})
	return Req{
		Params:  convertToParams(req["params"]),
		Headers: convertToHeaders(req["headers"]),
	}
}

func convertToRes(r interface{}) Res {
	res := r.(map[interface{}]interface{})
	return Res{
		Status:  res["status"].(int),
		Params:  convertToParams(res["params"]),
		Headers: convertToHeaders(res["headers"]),
	}
}

func convertToParams(p interface{}) map[string]string {
	params := make(map[string]string)
	if p != nil {
		for k, v := range p.(map[interface{}]interface{}) {
			params[k.(string)] = v.(string)
		}
	}
	return params
}

func convertToHeaders(h interface{}) map[string]string {
	headers := make(map[string]string)
	if h != nil {
		for k, v := range h.(map[interface{}]interface{}) {
			headers[k.(string)] = v.(string)
		}
	}
	return headers
}
