package atgen

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
	"strings"

	util "github.com/lkesteloot/astutil"
	"golang.org/x/tools/go/ast/astutil"
)

func (g *Generator) Generate() error {
	tfuncs := filterTestFuncs(g.TestFuncs)
	for v, t := range tfuncs {
		f, err := os.Create(fmt.Sprintf("%s_test.go", v))
		if err != nil {
			return err
		}
		defer f.Close()
		g.generateTestFuncs(v, t, f)
	}
	return nil
}

func (g *Generator) generateTestFuncs(version string, testFuncs TestFuncs, w io.Writer) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, g.Template, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	var (
		testFuncNode ast.Node
		testNode     ast.Node
	)

	cmap := ast.NewCommentMap(fset, f, f.Comments)
	for node, cgroups := range cmap {
		for _, cgroup := range cgroups {
			if strings.Contains(cgroup.Text(), "TestFunc block") {
				testFuncNode = node
			} else if strings.Contains(cgroup.Text(), "Test block") {
				testNode = node
			}
		}
	}

	astutil.Apply(testFuncNode, func(cr *astutil.Cursor) bool {
		if cr.Node() == testNode {
			cr.Delete()
		}
		return true
	}, nil)

	var tfnodes []ast.Node
	for _, testFunc := range testFuncs {
		tfnode := util.DuplicateNode(testFuncNode)
		tfnode.(*ast.FuncDecl).Name.Name = testFunc.Name

		var tnodes []ast.Node
		for _, test := range testFunc.Tests {
			var tnode ast.Node
			if !test.IsSubtests() {
				tnode = util.DuplicateNode(testNode)
				tnode = rewriteTestNode(tnode, test.(Test))
			}
			tnodes = append(tnodes, tnode)

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
		test := t.(Test)
		apiVersions := test.APIVersions
		test.Path = strings.Replace(test.Path, "{apiVersion}", version, 1)
		if contains(apiVersions, version) {
			tfunc.Tests = append(tfunc.Tests, test)
		} else if apiVersions == nil && contains(testFunc.APIVersions, version) {
			tfunc.Tests = append(tfunc.Tests, test)
		}
	}
	return tfunc
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
		_ = test
		versions = append(versions, test.(Test).APIVersions...)
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

func rewriteTestNode(n ast.Node, test Test) ast.Node {
	var ident string
	astutil.Apply(n, func(cr *astutil.Cursor) bool {
		switch v := cr.Node().(type) {
		case *ast.BasicLit:
			switch v.Value {
			case `"Method"`:
				v.Value = fmt.Sprintf(`"%s"`, strings.ToUpper(test.Method))
			case `"Path"`:
				v.Value = fmt.Sprintf(`"%s"`, test.Path)
			case "`reqParams`":
				params, _ := json.Marshal(test.Req.Params)
				v.Value = fmt.Sprintf("`%s`", params)
			}
		case *ast.Ident:
			ident = v.Name
		case *ast.CompositeLit:
			if ident == "reqHeaders" {
				h, _ := parser.ParseExpr(fmt.Sprintf("%#v", test.Req.Headers))
				cr.Replace(h)
			} else if ident == "resHeaders" {
				h, _ := parser.ParseExpr(fmt.Sprintf("%#v", test.Res.Headers))
				cr.Replace(h)
			} else if ident == "resParams" {
				p, _ := parser.ParseExpr(fmt.Sprintf("%#v", test.Res.Params))
				cr.Replace(p)
			}
			ident = ""
		}
		return true
	}, nil)
	return n
}
