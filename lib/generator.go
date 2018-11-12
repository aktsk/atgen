package atgen

import (
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

func (t TestFuncs) Generate(template string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, template, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	printer.Fprint(os.Stdout, fset, f)

	return nil
}
