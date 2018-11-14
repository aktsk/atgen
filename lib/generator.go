package atgen

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strings"

	util "github.com/lkesteloot/astutil"
	"golang.org/x/tools/go/ast/astutil"
)

func (g *Generator) Generate() error {
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
	for _, testFunc := range g.TestFuncs {
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

		astutil.Apply(tfnode, func(cr *astutil.Cursor) bool {
			switch v := cr.Node().(type) {
			case *ast.BlockStmt:
				if v.List == nil {
					for _, n := range tnodes {
						cr.InsertBefore(n)
					}
					cr.Delete()
				}
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
	printer.Fprint(os.Stdout, fset, f)

	return nil
}

func rewriteTestNode(n ast.Node, test Test) ast.Node {
	rewriteReqHeaders := false
	rewriteResHeaders := false
	rewriteResParams := false
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
			if v.Name == "reqHeaders" {
				rewriteReqHeaders = true
				rewriteResHeaders = false
				rewriteResParams = false
			} else if v.Name == "resHeaders" {
				rewriteReqHeaders = false
				rewriteResHeaders = true
				rewriteResParams = false
			} else if v.Name == "resParams" {
				rewriteReqHeaders = false
				rewriteResHeaders = false
				rewriteResParams = true
			}
		case *ast.CompositeLit:
			if rewriteReqHeaders {
				h, _ := parser.ParseExpr(fmt.Sprintf("%#v\n", test.Req.Headers))
				cr.Replace(h)
			} else if rewriteResHeaders {
				h, _ := parser.ParseExpr(fmt.Sprintf("%#v\n", test.Res.Headers))
				cr.Replace(h)
			} else if rewriteResParams {
				p, _ := parser.ParseExpr(fmt.Sprintf("%#v\n", test.Res.Params))
				cr.Replace(p)
			}
		}
		//fmt.Printf("%#v\n", cr.Node())
		return true
	}, nil)
	return n
}
