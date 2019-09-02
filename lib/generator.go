package atgen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	util "github.com/lkesteloot/astutil"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

// RouterFuncName is function name to be replaced
const RouterFuncName = "AtgenRouterFunc"

// Generate generates code and write to files
func (g *Generator) Generate() error {
	base := getFileNameWithoutExt(g.Yaml)
	if !strings.HasSuffix(base, "_test") {
		base = base + "_test"
	}

	tfuncs := filterTestFuncs(g.TestFuncs)
	for v, t := range tfuncs {
		filename := fmt.Sprintf("%s_%s.go", v, base)
		tf, err := ioutil.TempFile(g.OutputDir, filename)
		if err != nil {
			return errors.WithStack(err)
		}
		defer func() {
			tf.Close()
			os.Remove(tf.Name())
		}()
		err = g.generateTestFuncs(v, t, tf)
		if err != nil {
			return errors.WithStack(err)
		}
		out := filepath.Join(g.OutputDir, filename)
		f, err := os.Create(out)
		if err != nil {
			return errors.WithStack(err)
		}
		defer f.Close()

		tf.Seek(0, 0)
		io.Copy(f, tf)
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
			if strings.Contains(cgroup.Text(), "Atgen TestFunc block") {
				testFuncNode = node
			} else if strings.Contains(cgroup.Text(), "Atgen Test block") {
				testNode = node
			} else if strings.Contains(cgroup.Text(), "Atgen Subtest block") {
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

	absPath, err := filepath.Abs(g.OutputDir)
	if err != nil {
		return err
	}

	outputPath, err := PackageName(afero.NewOsFs(), os.Getenv("GOPATH"), absPath)
	if err != nil {
		return err
	}

	rewriteFileAst(fset, f, testFuncs, outputPath)

	var tfnodes []ast.Node
	for _, testFunc := range testFuncs {
		tfnode := util.DuplicateNode(testFuncNode)
		rewriteTestFuncNode(tfnode, testFunc, outputPath, g.Program)

		var tnodes []ast.Node
		for _, t := range testFunc.Tests {
			switch test := t.(type) {
			case Test:
				tnode := util.DuplicateNode(testNode)
				tnode = rewriteTestNode(tnode, test)
				tnodes = append(tnodes, tnode)
			case Subtests:
				for _, subtest := range test {
					subtnode := util.DuplicateNode(subtestNode)
					astutil.Apply(subtnode, func(cr *astutil.Cursor) bool {
						switch v := cr.Node().(type) {
						case *ast.BasicLit:
							switch v.Value {
							case `"AtgenSubtestName"`:
								v.Value = fmt.Sprintf(`"%s"`, subtest.Name)
							}
						}
						return true
					}, nil)

					var tests []ast.Node
					for _, test := range subtest.Tests {
						tnode := util.DuplicateNode(testNode)
						tnode = rewriteTestNode(tnode, test)
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
		Name:           testFunc.Name,
		Vars:           testFunc.Vars,
		RouterFuncName: testFunc.RouterFuncName,
		RouterFunc:     testFunc.RouterFunc,
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

func rewriteFileAst(fset *token.FileSet, f *ast.File, tfuncs TestFuncs, outputPath string) {
	for _, tfunc := range tfuncs {
		if tfunc.RouterFunc.PackagePath == outputPath {
			continue
		}
		// TODO: When package names conflict, this field should be set with a generated unique name
		astutil.AddImport(fset, f, tfunc.RouterFunc.PackagePath)
	}
}

func rewriteTestFuncNode(n ast.Node, tfunc TestFunc, outputPath string, pkgs []*packages.Package) {
	n.(*ast.FuncDecl).Name.Name = tfunc.Name
	astutil.Apply(n, func(cr *astutil.Cursor) bool {
		switch v := cr.Node().(type) {
		case *ast.CallExpr:
			ident, ok := v.Fun.(*ast.Ident)
			if ok && ident.Name == RouterFuncName {
				if tfunc.RouterFunc.PackagePath == outputPath {
					v.Fun = &ast.Ident{Name: tfunc.RouterFunc.Name}
				} else {
					var pkg *packages.Package
					for _, p := range pkgs {
						if p.PkgPath == tfunc.RouterFunc.PackagePath {
							pkg = p
						}
					}
					v.Fun = &ast.SelectorExpr{
						X:   &ast.Ident{Name: pkg.Name},
						Sel: &ast.Ident{Name: tfunc.RouterFunc.Name},
					}
				}
			}
		}
		return true
	}, nil)
}

func rewriteTestNode(n ast.Node, test Test) ast.Node {
	var ident string
	astutil.Apply(n, func(cr *astutil.Cursor) bool {
		switch v := cr.Node().(type) {
		case *ast.BasicLit:
			switch v.Value {
			case `"AtgenMethod"`:
				v.Value = fmt.Sprintf(`"%s"`, strings.ToUpper(test.Method))
			case `"AtgenPath"`:
				v.Value = fmt.Sprintf(`"%s"`, test.Path)
			case `"atgenStatus"`:
				v.Value = fmt.Sprintf("%d", test.Res.Status)
			case `"atgenRegisterKey"`:
				v.Value = fmt.Sprintf(`"%s"`, test.Register)
			}
		case *ast.Ident:
			ident = v.Name
		case *ast.CompositeLit:
			if ident == "atgenReqHeaders" {
				h, _ := parser.ParseExpr(fmt.Sprintf("%#v", test.Req.Headers))
				cr.Replace(h)
			} else if ident == "atgenReqParams" {
				p, _ := parser.ParseExpr(fmt.Sprintf("%#v", test.Req.Params))
				cr.Replace(p)
			} else if ident == "atgenResHeaders" {
				h, _ := parser.ParseExpr(fmt.Sprintf("%#v", test.Res.Headers))
				cr.Replace(h)
			} else if ident == "atgenResParams" {
				p, _ := parser.ParseExpr(fmt.Sprintf("%#v", test.Res.Params))
				cr.Replace(p)
			} else if ident == "atgenTestVars" {
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
				s := strings.TrimPrefix(v.Value, `"${`)
				s = strings.TrimSuffix(s, `}"`)
				t := strings.Split(s, ":")
				v.Value = fmt.Sprintf(`vars["%s"].(%s)`, t[0], t[1])
			} else if strings.HasPrefix(v.Value, `"$atgenRegister[`) {
				s := strings.TrimPrefix(v.Value, `"$atgenRegister[`)
				s = strings.TrimSuffix(s, `]"`)
				t := strings.Split(s, ".")
				v.Value = fmt.Sprintf(`atgenRegister["%s"].(map[string]interface{})["%s"].(string)`, t[0], t[1])
			}
		}
		return true
	}, nil)

	return n
}
