package atgen

import (
	"io/ioutil"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// ParseYaml parses yaml which deifnes test requests/responses
// and convert it to types defined in types.go
func (g *Generator) ParseYaml() error {
	buf, err := ioutil.ReadFile(g.Yaml)
	if err != nil {
		return errors.WithStack(err)
	}

	parsed, err := parseYaml(buf)
	if err != nil {
		return errors.WithStack(err)
	}
	testFuncs, err := convertToTestFuncs(parsed)
	if err != nil {
		return errors.WithStack(err)
	}

	g.TestFuncs = testFuncs

	return nil
}

func parseYaml(buf []byte) ([]map[interface{}]interface{}, error) {
	var parsed []map[interface{}]interface{}
	err := yaml.Unmarshal(buf, &parsed)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return parsed, nil
}

func convertToTestFuncs(parsed []map[interface{}]interface{}) (TestFuncs, error) {
	var testFuncs TestFuncs
	for _, p := range parsed {
		if p["name"] == nil {
			continue
		}

		testFunc := TestFunc{}
		name, ok := p["name"].(string)
		if !ok {
			return testFuncs, errors.New("name must be string")
		}
		testFunc.Name = name

		routerFunc, ok := p["routerFunc"].(string)
		if !ok {
			return testFuncs, errors.New("routerFunc must be string")
		}
		testFunc.RouterFunc = routerFunc

		if p["apiVersions"] != nil {
			for _, v := range p["apiVersions"].([]interface{}) {
				testFunc.APIVersions = append(testFunc.APIVersions, v.(string))
			}
		}

		vars, err := convertToParams(p["vars"])
		if err != nil {
			return testFuncs, err
		}
		testFunc.Vars = vars

		for _, t := range p["tests"].([]interface{}) {
			t := t.(map[interface{}]interface{})
			if t["path"] != nil {
				test, err := convertToTest(t)
				if err != nil {
					return testFuncs, err
				}
				testFunc.Tests = append(testFunc.Tests, test)
			} else {
				subtests, err := convertToSubtests(t)
				if err != nil {
					return testFuncs, err
				}
				testFunc.Tests = append(testFunc.Tests, subtests)
			}
		}
		testFuncs = append(testFuncs, testFunc)
	}
	return testFuncs, nil
}

func convertToTest(t map[interface{}]interface{}) (Test, error) {
	var apiVersions []string
	if t["apiVersions"] != nil {
		for _, v := range t["apiVersions"].([]interface{}) {
			apiVersions = append(apiVersions, v.(string))
		}
	}

	req, err := convertToReq(t["req"])
	if err != nil {
		return Test{}, err
	}
	res, err := convertToRes(t["res"])
	if err != nil {
		return Test{}, err
	}
	vars, err := convertToParams(t["vars"])
	if err != nil {
		return Test{}, err
	}

	return Test{
		APIVersions: apiVersions,
		Path:        t["path"].(string),
		Method:      t["method"].(string),
		Req:         req,
		Res:         res,
		Vars:        vars,
	}, nil
}

func convertToSubtests(s map[interface{}]interface{}) (Subtests, error) {
	subTests := Subtests{}

	for _, s := range s["subtests"].([]interface{}) {
		s := s.(map[interface{}]interface{})
		subtest := Subtest{Name: s["name"].(string)}
		for _, t := range s["tests"].([]interface{}) {
			t := t.(map[interface{}]interface{})
			test, err := convertToTest(t)
			if err != nil {
				return subTests, err
			}
			subtest.Tests = append(subtest.Tests, test)
		}
		subTests = append(subTests, subtest)
	}
	return subTests, nil
}

func convertToReq(r interface{}) (Req, error) {
	if r == nil {
		return Req{}, nil
	}

	req := r.(map[interface{}]interface{})

	headers, err := convertToHeaders(req["headers"])
	if err != nil {
		return Req{}, nil
	}

	params, err := convertToParams(req["params"])
	if err != nil {
		return Req{}, err
	}

	return Req{
		Params:  params,
		Headers: headers,
	}, nil
}

func convertToRes(r interface{}) (Res, error) {
	if r == nil {
		return Res{}, nil
	}

	res := r.(map[interface{}]interface{})

	headers, err := convertToHeaders(res["headers"])
	if err != nil {
		return Res{}, err
	}

	params, err := convertToParams(res["params"])
	if err != nil {
		return Res{}, err
	}

	return Res{
		Status:  res["status"].(int),
		Params:  params,
		Headers: headers,
	}, nil
}

func convertToParams(p interface{}) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	if p != nil {
		for k, v := range p.(map[interface{}]interface{}) {
			key, ok := k.(string)
			if !ok {
				return params, errors.New("key should be string")
			}

			switch t := v.(type) {
			case string, bool, int:
				params[key] = t
			case map[interface{}]interface{}:
				p, err := convertToParams(t)
				if err != nil {
					return params, err
				}
				params[key] = p
			case []interface{}:
				params[key] = []map[string]interface{}{}
				for _, v := range t {
					p, err := convertToParams(v)
					if err != nil {
						return params, err
					}
					params[key] = append(params[key].([]map[string]interface{}), p)
				}
			default:
				return params, errors.New("invalid type")
			}
		}
	}
	return params, nil
}

func convertToHeaders(h interface{}) (map[string]string, error) {
	headers := make(map[string]string)
	if h != nil {
		for k, v := range h.(map[interface{}]interface{}) {
			key, ok := k.(string)
			if !ok {
				return headers, errors.New("header key must be string")
			}

			val, ok := v.(string)
			if !ok {
				return headers, errors.New("header val must be string")
			}
			headers[key] = val
		}
	}
	return headers, nil
}
