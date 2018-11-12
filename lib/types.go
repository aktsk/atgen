package atgen

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
	Params  Params
	Headers Headers
}

type Res struct {
	Status  int
	Params  Params
	Headers Headers
}

type Params map[string]string
type Headers map[string]string

type Tester interface {
	IsSubtests() bool
}

func (t Test) IsSubtests() bool {
	return false
}

func (t SubTests) IsSubtests() bool {
	return true
}
