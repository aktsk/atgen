package atgen

import (
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
		if cr.Node() == subtestNode {
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
	astutil.Apply(n, func(cr *astutil.Cursor) bool {
		// TODO: rewrite test node
		return true
	}, nil)
	return n
}
