package atgen

import (
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"golang.org/x/tools/go/packages"
	yaml "gopkg.in/yaml.v3"
)

// ParseYaml parses yaml which defines test requests/responses
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

	absPath, err := filepath.Abs(g.TemplateDir)
	if err != nil {
		return errors.WithStack(err)
	}

	packageName, err := PackageName(afero.NewOsFs(), os.Getenv("GOPATH"), absPath)
	if err != nil {
		return errors.WithStack(err)
	}

	routerFuncs, testFuncs, err := aggregateRouterFunc(testFuncs, packageName)
	if err != nil {
		return errors.WithStack(err)
	}

	program, err := loadUsingPackages(routerFuncs)
	if err != nil {
		return errors.WithStack(err)
	}

	err = validateRouterFuncs(routerFuncs, program)
	if err != nil {
		return errors.WithStack(err)
	}

	g.Program = program
	g.TestFuncs = testFuncs

	return nil
}

func parseYaml(buf []byte) ([]map[string]interface{}, error) {
	var parsed []map[string]interface{}
	err := yaml.Unmarshal(buf, &parsed)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return parsed, nil
}

func convertToTestFuncs(parsed []map[string]interface{}) (TestFuncs, error) {
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
		testFunc.RouterFuncName = routerFunc

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
			t := t.(map[string]interface{})
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

func loadUsingPackages(routerFuncs []*RouterFunc) ([]*packages.Package, error) {
	packagePaths := []string{"net/http"}

	for _, routerFunc := range routerFuncs {
		if !usingSamePackage(packagePaths, routerFunc) {
			packagePaths = append(packagePaths, routerFunc.PackagePath)
		}
	}

	conf := &packages.Config{}
	pkgs, err := packages.Load(conf, packagePaths...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, pkg := range pkgs {
		if pkg.Errors != nil {
			return nil, errors.WithStack(pkg.Errors[0])
		}
	}

	return pkgs, nil
}

func usingSamePackage(packagePaths []string, routerFunc *RouterFunc) bool {
	for _, packagePath := range packagePaths {
		if packagePath == routerFunc.PackagePath {
			return true
		}
	}
	return false
}

func aggregateRouterFunc(tfuncs TestFuncs, templatePath string) ([]*RouterFunc, TestFuncs, error) {
	newTfuncs := make(TestFuncs, len(tfuncs))
	routerFuncs := []*RouterFunc{}
	for i, tfunc := range tfuncs {
		routerFuncNameList := strings.Split(tfunc.RouterFuncName, ".")
		if len(routerFuncNameList) == 0 {
			return nil, nil, errors.New("invalid format routerFunc")
		}
		var packagePath, funcName string
		if len(routerFuncNameList) == 1 {
			packagePath = "./"
			funcName = routerFuncNameList[0]
		} else {
			packagePath = strings.Join(routerFuncNameList[0:len(routerFuncNameList)-1], ".")
			funcName = routerFuncNameList[len(routerFuncNameList)-1]
		}

		if isRelativePath(packagePath) {
			packagePath = filepath.Join(templatePath, packagePath)
		}

		routerFunc := RouterFunc{
			PackagePath: packagePath,
			Name:        funcName,
		}
		tfunc.RouterFunc = &routerFunc
		newTfuncs[i] = tfunc
		routerFuncs = append(routerFuncs, &routerFunc)
	}

	return routerFuncs, newTfuncs, nil
}

func isRelativePath(path string) bool {
	return strings.HasPrefix(path, ".")
}

var packageCache = map[string][]*packages.Package{}

func validateRouterFuncs(routerFuncs []*RouterFunc, program []*packages.Package) error {
	conf := &packages.Config{Mode: packages.LoadAllSyntax}
	for _, routerFunc := range routerFuncs {
		path := routerFunc.PackagePath
		if packageCache[path] == nil {
			pkgs, err := packages.Load(conf, "net/http", path)
			if err != nil {
				return errors.WithStack(err)
			}
			packageCache[path] = pkgs
		}

		pkgs := packageCache[path]
		handlerObj := pkgs[0].Types.Scope().Lookup("Handler")
		funcObj := pkgs[1].Types.Scope().Lookup(routerFunc.Name)

		err := validateRouterFuncObj(handlerObj, funcObj, routerFunc)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateRouterFuncObj(handlerObj types.Object, routerFuncObj types.Object, routerFunc *RouterFunc) error {
	if routerFuncObj == nil {
		return errors.Errorf("can't resolve %s.%s", routerFunc.PackagePath, routerFunc.Name)
	}

	fun, ok := routerFuncObj.(*types.Func)
	if !ok {
		return errors.Errorf("%s.%s is not function", routerFunc.PackagePath, routerFunc.Name)
	}

	signature, ok := fun.Type().(*types.Signature)
	if !ok {
		return errors.Errorf("%s.%s is not signature", routerFunc.PackagePath, routerFunc.Name)
	}

	if signature.Recv() != nil {
		return errors.Errorf("%s.%s signature should be func() http.Handler (the function is method)", routerFunc.PackagePath, routerFunc.Name)
	}

	if signature.Params().Len() != 0 {
		return errors.Errorf("%s.%s signature should be func() http.Handler (the function receive params)", routerFunc.PackagePath, routerFunc.Name)
	}

	if signature.Results().Len() != 1 {
		return errors.Errorf("%s.%s signature should be func() http.Handler (the function return more than 1 results)", routerFunc.PackagePath, routerFunc.Name)
	}

	if !types.AssignableTo(signature.Results().At(0).Type(), handlerObj.Type()) {
		return errors.Errorf("%s.%s signature should be func() http.Handler (the function don't return http.Handler)", routerFunc.PackagePath, routerFunc.Name)
	}

	return nil
}

func convertToTest(t map[string]interface{}) (Test, error) {
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

	register := ""
	if t["register"] != nil {
		register = t["register"].(string)
	}

	return Test{
		APIVersions: apiVersions,
		Path:        t["path"].(string),
		Method:      t["method"].(string),
		Req:         req,
		Res:         res,
		Vars:        vars,
		Register:    register,
	}, nil
}

func convertToSubtests(s map[string]interface{}) (Subtests, error) {
	subTests := Subtests{}

	for _, s := range s["subtests"].([]interface{}) {
		s := s.(map[string]interface{})
		subtest := Subtest{Name: s["name"].(string)}
		for _, t := range s["tests"].([]interface{}) {
			t := t.(map[string]interface{})
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

	req := r.(map[string]interface{})

	headers, err := convertToHeaders(req["headers"])
	if err != nil {
		return Req{}, nil
	}

	params, err := convertToParams(req["params"])
	if err != nil {
		return Req{}, err
	}

	body := ""
	if b, ok := req["body"].(string); ok {
		body = b
	}

	typ := JSON
	contentType := ""
	if typStr, ok := req["type"].(string); ok {
		switch typStr {
		case "json":
			typ = JSON
			contentType = "application/json"
		case "form":
			typ = FORM
			contentType = "application/x-www-form-urlencoded"
		case "raw":
			typ = RAW
			contentType = "application/octet-stream"
		default:
			return Req{}, errors.New("request type must be json, form or raw")
		}
	} else {
		return Req{}, errors.New("request type must be defined")
	}

	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = contentType
	}

	return Req{
		Params:  params,
		Headers: headers,
		Body:    body,
		Type:    typ,
	}, nil
}

func convertToRes(r interface{}) (Res, error) {
	if r == nil {
		return Res{}, nil
	}

	res := r.(map[string]interface{})

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
		for key, v := range p.(map[string]interface{}) {
			switch t := v.(type) {
			case string, bool, int:
				params[key] = t
			case map[string]interface{}:
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
		for key, v := range h.(map[string]interface{}) {
			val, ok := v.(string)
			if !ok {
				return headers, errors.New("header val must be string")
			}
			headers[key] = val
		}
	}
	return headers, nil
}
