package atgen

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strings"

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
		//		testNode     ast.Node
		//		subTestsNode ast.Node
	)

	cmap := ast.NewCommentMap(fset, f, f.Comments)
	for node, cgroups := range cmap {
		for _, cgroup := range cgroups {
			if strings.Contains(cgroup.Text(), "Begin TestFunc") {
				testFuncNode = node
			}
		}
	}

	var nodes []ast.Node
	for _, testFunc := range g.TestFuncs {
		n := testFuncNode
		n.(*ast.FuncDecl).Name.Name = testFunc.Name
		for _, test := range testFunc.Tests {
			_ = test
			//fmt.Printf("%#v\n", test)
			//fmt.Println(test.IsSubtests())
		}
		nodes = append(nodes, n)
	}

	n := astutil.Apply(f, func(cr *astutil.Cursor) bool {
		if cr.Name() == "Decls" {
			switch v := cr.Node().(type) {
			case *ast.GenDecl:
				if v.Tok == token.IMPORT {
					for _, n := range nodes {
						cr.InsertAfter(n)
					}
				}
			}
		}
		return true
	}, nil)

	printer.Fprint(os.Stdout, fset, n)

	return nil
}
