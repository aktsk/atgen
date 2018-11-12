package atgen

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

func ParseYaml(file string) ([]TestFunc, error) {
	var testFuncs []TestFunc

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var parsed []map[interface{}]interface{}
	err = yaml.Unmarshal(buf, &parsed)
	if err != nil {
		return nil, err
	}

	for _, p := range parsed {
		testFunc := TestFunc{}
		testFunc.Name = p["name"].(string)

		for _, t := range p["tests"].([]interface{}) {
			t := t.(map[interface{}]interface{})
			if t["path"] != nil {
				testFunc.Tests = append(testFunc.Tests, convertToTest(t))
			}
		}

		testFuncs = append(testFuncs, testFunc)
	}

	return testFuncs, nil
}

func convertToTest(t map[interface{}]interface{}) Test {
	return Test{
		Path:   t["path"].(string),
		Method: t["method"].(string),
		Req:    convertToReq(t["req"]),
		Res:    convertToRes(t["res"]),
	}
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

func convertToParams(p interface{}) Params {
	params := Params{}
	if p != nil {
		for k, v := range p.(map[interface{}]interface{}) {
			params[k.(string)] = v.(string)
		}
	}
	return params
}

func convertToHeaders(h interface{}) Headers {
	headers := Headers{}
	if h != nil {
		for k, v := range h.(map[interface{}]interface{}) {
			headers[k.(string)] = v.(string)
		}
	}
	return headers
}
