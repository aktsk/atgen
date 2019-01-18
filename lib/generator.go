package atgen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"

	util "github.com/lkesteloot/astutil"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/ast/astutil"
)

const RouterFuncName = "routerFunc"

// Generate generates code and write to files
func (g *Generator) Generate() error {
	base := getFileNameWithoutExt(g.Yaml)
	if !strings.HasSuffix(base, "_test") {
		base = base + "_test"
	}

	tfuncs := filterTestFuncs(g.TestFuncs)
	for v, t := range tfuncs {
		out := filepath.Join(g.OutputDir, fmt.Sprintf("%s_%s.go", v, base))
		f, err := os.Create(out)
		if err != nil {
			return errors.WithStack(err)
		}
		defer f.Close()
		err = g.generateTestFuncs(v, t, f)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func getFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

func (g *Generator) generateTestFuncs(version string, testFuncs TestFuncs, w io.Writer) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, g.Template, nil, parser.ParseComments)
	if err != nil {
		return errors.WithStack(err)
	}

	var (
		testFuncNode ast.Node
		testNode     ast.Node
		subtestNode  ast.Node
	)

	cmap := ast.NewCommentMap(fset, f, f.Comments)
	for node, cgroups := range cmap {
		for _, cgroup := range cgroups {
			if strings.Contains(cgroup.Text(), "TestFunc block") {
				testFuncNode = node
			} else if strings.Contains(cgroup.Text(), "Test block") {
				testNode = node
			} else if strings.Contains(cgroup.Text(), "Subtest block") {
				subtestNode = node
			}
		}
	}

	astutil.Apply(testFuncNode, func(cr *astutil.Cursor) bool {

		if cr.Node() == testNode {
			cr.Delete()
		}

		if cr.Node() == subtestNode && subtestNode != nil {
			cr.Delete()
		}

		return true
	}, nil)

	var tfnodes []ast.Node
	for _, testFunc := range testFuncs {
		tfnode := util.DuplicateNode(testFuncNode)
		tfnode.(*ast.FuncDecl).Name.Name = testFunc.Name

		var tnodes []ast.Node
		for _, t := range testFunc.Tests {
			switch test := t.(type) {
			case Test:
				tnode := util.DuplicateNode(testNode)
				tnode = rewriteTestNode(tnode, test, testFunc)
				tnodes = append(tnodes, tnode)
			case Subtests:
				for _, subtest := range test {
					subtnode := util.DuplicateNode(subtestNode)
					astutil.Apply(subtnode, func(cr *astutil.Cursor) bool {
						switch v := cr.Node().(type) {
						case *ast.BasicLit:
							switch v.Value {
							case `"SubtestName"`:
								v.Value = fmt.Sprintf(`"%s"`, subtest.Name)
							}
						}
						return true
					}, nil)

					var tests []ast.Node
					for _, test := range subtest.Tests {
						tnode := util.DuplicateNode(testNode)
						tnode = rewriteTestNode(tnode, test, testFunc)
						tests = append(tests, tnode)
					}
					subtnode = rewriteSubtestNode(subtnode, tests)
					tnodes = append(tnodes, subtnode)
				}
			}
		}

		var ident string
		astutil.Apply(tfnode, func(cr *astutil.Cursor) bool {
			switch v := cr.Node().(type) {
			case *ast.BlockStmt:
				if v.List == nil {
					for _, n := range tnodes {
						cr.InsertBefore(n)
					}
					cr.Delete()
				}
			case *ast.Ident:
				ident = v.Name
			case *ast.CompositeLit:
				if ident == "vars" {
					h, _ := parser.ParseExpr(fmt.Sprintf("%#v", testFunc.Vars))
					cr.Replace(h)
				}
				ident = ""
			}

			return true
		}, nil)

		tfnodes = append(tfnodes, tfnode)
	}

	astutil.Apply(f, func(cr *astutil.Cursor) bool {
		if cr.Name() == "Decls" {
			switch cr.Node().(type) {
			case *ast.FuncDecl:
				for _, n := range tfnodes {
					cr.InsertBefore(n)
				}
			}
		}

		if cr.Node() == testFuncNode {
			cr.Delete()
		}

		return true
	}, nil)

	f.Comments = cmap.Filter(f).Comments()
	printer.Fprint(w, fset, f)

	return nil
}

func rewriteSubtestNode(subtest ast.Node, tests []ast.Node) ast.Node {
	astutil.Apply(subtest, func(cr *astutil.Cursor) bool {
		switch v := cr.Node().(type) {
		case *ast.BlockStmt:
			if v.List == nil {
				for _, n := range tests {
					cr.InsertBefore(n)
				}
				cr.Delete()
			}
		}
		return true
	}, nil)

	return subtest
}

func filterTestFuncs(testFuncs TestFuncs) map[string]TestFuncs {
	tfuncs := make(map[string]TestFuncs)
	for _, testFunc := range testFuncs {
		for _, version := range getVersions(testFunc) {
			tfunc := filterTests(testFunc, version)
			tfuncs[version] = append(tfuncs[version], tfunc)
		}
	}
	return tfuncs
}

func filterTests(testFunc TestFunc, version string) TestFunc {
	tfunc := TestFunc{
		Name: testFunc.Name,
		Vars: testFunc.Vars,
	}
	for _, t := range testFunc.Tests {
		switch v := t.(type) {
		case Test:
			test := filterTest(v, testFunc.APIVersions, version)
			if test != nil {
				tfunc.Tests = append(tfunc.Tests, *test)
			}
		case Subtests:
			subtests := Subtests{}
			for _, s := range v {
				subtest := Subtest{Name: s.Name}
				if s.APIVersions != nil && !contains(s.APIVersions, version) {
					continue
				}
				if s.APIVersions == nil && !contains(testFunc.APIVersions, version) {
					continue
				}
				for _, t := range s.Tests {
					test := filterTest(t, testFunc.APIVersions, version)
					if test != nil {
						subtest.Tests = append(subtest.Tests, *test)
					}
				}
				subtests = append(subtests, subtest)
			}
			tfunc.Tests = append(tfunc.Tests, subtests)
		}
	}
	return tfunc
}

func filterTest(test Test, versions []string, version string) *Test {
	apiVersions := test.APIVersions
	test.Path = strings.Replace(test.Path, "{apiVersion}", version, 1)
	if contains(apiVersions, version) {
		return &test
	}

	if apiVersions == nil && contains(versions, version) {
		return &test
	}

	return nil
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func getVersions(testFunc TestFunc) []string {
	var versions []string
	versions = append(versions, testFunc.APIVersions...)
	for _, test := range testFunc.Tests {
		switch v := test.(type) {
		case Test:
			versions = append(versions, v.APIVersions...)
		}
	}

	// Dedupe versions
	m := make(map[string]bool)
	var deduped []string
	for _, v := range versions {
		if !m[v] {
			m[v] = true
			deduped = append(deduped, v)
		}
	}

	return deduped
}

func rewriteTestNode(n ast.Node, test Test, tfunc TestFunc) ast.Node {
	var ident string
	astutil.Apply(n, func(cr *astutil.Cursor) bool {
		switch v := cr.Node().(type) {
		case *ast.CallExpr:
			ident, ok := v.Fun.(*ast.Ident)
			if ok && ident.Name == RouterFuncName {
				ident.Name = tfunc.RouterFunc
			}
		case *ast.BasicLit:
			switch v.Value {
			case `"Method"`:
				v.Value = fmt.Sprintf(`"%s"`, strings.ToUpper(test.Method))
			case `"Path"`:
				v.Value = fmt.Sprintf(`"%s"`, test.Path)
			case `"status"`:
				v.Value = fmt.Sprintf("%d", test.Res.Status)
			}
		case *ast.Ident:
			ident = v.Name
		case *ast.CompositeLit:
			if ident == "reqHeaders" {
				h, _ := parser.ParseExpr(fmt.Sprintf("%#v", test.Req.Headers))
				cr.Replace(h)
			} else if ident == "reqParams" {
				p, _ := parser.ParseExpr(fmt.Sprintf("%#v", test.Req.Params))
				cr.Replace(p)
			} else if ident == "resHeaders" {
				h, _ := parser.ParseExpr(fmt.Sprintf("%#v", test.Res.Headers))
				cr.Replace(h)
			} else if ident == "resParams" {
				p, _ := parser.ParseExpr(fmt.Sprintf("%#v", test.Res.Params))
				cr.Replace(p)
			} else if ident == "testVars" {
				h, _ := parser.ParseExpr(fmt.Sprintf("%#v", test.Vars))
				cr.Replace(h)
			}

			ident = ""
		}
		return true
	}, nil)

	astutil.Apply(n, func(cr *astutil.Cursor) bool {
		switch v := cr.Node().(type) {
		case *ast.BasicLit:
			if strings.HasPrefix(v.Value, `"${`) {
				s := strings.TrimLeft(v.Value, `"${`)
				s = strings.TrimRight(s, `}"`)
				t := strings.Split(s, ":")
				v.Value = fmt.Sprintf(`vars["%s"].(%s)`, t[0], t[1])
			}
		}
		return true
	}, nil)

	return n
}
