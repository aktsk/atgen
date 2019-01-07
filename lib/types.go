package atgen

type Generator struct {
	Yaml                   string
	Template               string
	TemplateDir            string
	OutputDir              string
	TestFuncs              TestFuncs
	TestFuncsPerAPIVersion map[string]TestFuncs
}

type TestFuncs []TestFunc

type TestFunc struct {
	Name        string
	Tests       []Tester
	APIVersions []string
	Vars        map[string]interface{}
}

type Test struct {
	APIVersions []string
	Path        string
	Method      string
	Req         Req
	Res         Res
	Vars        map[string]interface{}
}

type Subtests []Subtest

type Subtest struct {
	Name        string
	Tests       []Test
	APIVersions []string
}

type Req struct {
	Params  map[string]interface{}
	Headers map[string]string
}

type Res struct {
	Status  int
	Params  map[string]interface{}
	Headers map[string]string
}

type Tester interface {
	IsSubtests() bool
}

func (t Test) IsSubtests() bool {
	return false
}

func (t Subtests) IsSubtests() bool {
	return true
}
