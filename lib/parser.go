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
		if p["name"] == nil {
			continue
		}

		testFunc := TestFunc{}
		testFunc.Name = p["name"].(string)

		if p["apiVersions"] != nil {
			for _, v := range p["apiVersions"].([]interface{}) {
				testFunc.APIVersions = append(testFunc.APIVersions, v.(string))
			}
		}

		testFunc.Vars = convertToParams(p["vars"])

		for _, t := range p["tests"].([]interface{}) {
			t := t.(map[interface{}]interface{})
			if t["path"] != nil {
				testFunc.Tests = append(testFunc.Tests, convertToTest(t))
			} else {
				testFunc.Tests = append(testFunc.Tests, convertToSubtests(t))
			}
		}
		testFuncs = append(testFuncs, testFunc)
	}
	return testFuncs
}

func convertToTest(t map[interface{}]interface{}) Test {
	var apiVersions []string
	if t["apiVersions"] != nil {
		for _, v := range t["apiVersions"].([]interface{}) {
			apiVersions = append(apiVersions, v.(string))
		}
	}

	return Test{
		APIVersions: apiVersions,
		Path:        t["path"].(string),
		Method:      t["method"].(string),
		Req:         convertToReq(t["req"]),
		Res:         convertToRes(t["res"]),
	}
}

func convertToSubtests(s map[interface{}]interface{}) Subtests {
	subTests := Subtests{}

	for _, s := range s["subtests"].([]interface{}) {
		s := s.(map[interface{}]interface{})
		subtest := Subtest{Name: s["name"].(string)}
		for _, t := range s["tests"].([]interface{}) {
			t := t.(map[interface{}]interface{})
			subtest.Tests = append(subtest.Tests, convertToTest(t))
		}
		subTests = append(subTests, subtest)
	}
	return subTests
}

func convertToReq(r interface{}) Req {
	if r == nil {
		return Req{}
	}

	req := r.(map[interface{}]interface{})
	return Req{
		Params:  convertToParams(req["params"]),
		Headers: convertToHeaders(req["headers"]),
	}
}

func convertToRes(r interface{}) Res {
	if r == nil {
		return Res{}
	}

	res := r.(map[interface{}]interface{})
	return Res{
		Status:  res["status"].(int),
		Params:  convertToParams(res["params"]),
		Headers: convertToHeaders(res["headers"]),
	}
}

func convertToParams(p interface{}) map[string]interface{} {
	params := make(map[string]interface{})
	if p != nil {
		for k, v := range p.(map[interface{}]interface{}) {
			switch t := v.(type) {
			case string, bool, int:
				params[k.(string)] = t
			case map[interface{}]interface{}:
				params[k.(string)] = convertToParams(t)
			case []interface{}:
				params[k.(string)] = []map[string]interface{}{}
				for _, v := range t {
					params[k.(string)] = append(params[k.(string)].([]map[string]interface{}), convertToParams(v))
				}
			}
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
