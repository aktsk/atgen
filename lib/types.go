package atgen

type Generator struct {
	Yaml      string
	Template  string
	TestFuncs TestFuncs
}

type TestFuncs []TestFunc

type TestFunc struct {
	Name  string
	Tests []Tester
}

type Test struct {
	Path   string
	Method string
	Req    Req
	Res    Res
}

type SubTests []Test

type Req struct {
	Params  map[string]string
	Headers map[string]string
}

type Res struct {
	Status  int
	Params  map[string]string
	Headers map[string]string
}

type Tester interface {
	IsSubtests() bool
}

func (t Test) IsSubtests() bool {
	return false
}

func (t SubTests) IsSubtests() bool {
	return true
}
