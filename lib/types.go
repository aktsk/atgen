package atgen

import "golang.org/x/tools/go/loader"

// Generator is the type for code generator
type Generator struct {
	Yaml                   string
	Template               string
	TemplateDir            string
	OutputDir              string
	RouterFuncs            []*RouterFunc
	TestFuncs              TestFuncs
	TestFuncsPerAPIVersion map[string]TestFuncs
	Program                *loader.Program
}

// TestFuncs is a group of TestFunc
type TestFuncs []TestFunc

// TestFunc represents a test function
type TestFunc struct {
	Name           string
	Tests          []Tester
	APIVersions    []string
	RouterFuncName string
	RouterFunc     *RouterFunc
	Vars           map[string]interface{}
}

// Test represents a test in a test function
type Test struct {
	APIVersions []string
	Path        string
	Method      string
	Req         Req
	Res         Res
	Vars        map[string]interface{}
}

// Subtests is a group of Subtest
type Subtests []Subtest

// Subtest reppresents a subtest
type Subtest struct {
	Name        string
	Tests       []Test
	APIVersions []string
}

// Req is a request parameters and headers which a test should throw
type Req struct {
	Params  map[string]interface{}
	Headers map[string]string
}

// Res is a response status, parameters and headers which a test should get
type Res struct {
	Status  int
	Params  map[string]interface{}
	Headers map[string]string
}

// Tester is an interface for Test and Subtest
type Tester interface {
	IsSubtests() bool
}

// IsSubtests returns false when t is Test
func (t Test) IsSubtests() bool {
	return false
}

// IsSubtests returns true when t is Subtests
func (t Subtests) IsSubtests() bool {
	return true
}

// RouterFunc describe a function which should be called from test to get http.Handler
type RouterFunc struct {
	PackagePath string
	Name        string
}
