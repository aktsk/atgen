package atgen

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

func (g *Generator) ParseYaml() error {
	var testFuncs TestFuncs

	buf, err := ioutil.ReadFile(g.Yaml)
	if err != nil {
		return err
	}

	var parsed []map[interface{}]interface{}
	err = yaml.Unmarshal(buf, &parsed)
	if err != nil {
		return err
	}

	for _, p := range parsed {
		testFunc := TestFunc{}
		testFunc.Name = p["name"].(string)
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

	g.TestFuncs = testFuncs

	return nil
}

func convertToTest(t map[interface{}]interface{}) Test {
	return Test{
		Path:   t["path"].(string),
		Method: t["method"].(string),
		Req:    convertToReq(t["req"]),
		Res:    convertToRes(t["res"]),
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
