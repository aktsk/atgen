package atgen

import (
	"io"
	"os"
	"testing"

	"github.com/spf13/afero"
)

const testGoMod = `
module github.com/aktsk/atgen

require (
	github.com/aktsk/atgen v0.1.0 // indirect
	github.com/gorilla/mux v1.6.2 // indirect
	github.com/lkesteloot/astutil v0.0.0-20130122170032-b6715328cfa5
	github.com/pkg/errors v0.8.1
	github.com/spf13/afero v1.2.0 // indirect
	github.com/urfave/cli v1.20.0
	golang.org/x/tools v0.0.0-20190117194123-b4b6fe2cb829
	gopkg.in/yaml.v2 v2.2.2
)
`

func TestPackageName(t *testing.T) {
	testCases := []struct {
		Name         string
		Gopath       string
		FS           afero.Fs
		CurrentPath  string
		ExpectedPath string
	}{
		{
			Name:         "in GOPATH",
			Gopath:       "/go",
			FS:           gopathFS(),
			CurrentPath:  "/go/src/github.com/aktsk/atgen",
			ExpectedPath: "github.com/aktsk/atgen",
		},
		{
			Name:         "using go.mod",
			Gopath:       "/go",
			FS:           moduleFS(),
			CurrentPath:  "/atgen",
			ExpectedPath: "github.com/aktsk/atgen",
		},
		{
			Name:         "subdir of package of using go.mod",
			Gopath:       "/go",
			FS:           moduleFS(),
			CurrentPath:  "/atgen/lib",
			ExpectedPath: "github.com/aktsk/atgen/lib",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			path, err := PackageName(testCase.FS, testCase.Gopath, testCase.CurrentPath)
			if err != nil {
				t.Fatalf("%+v", err)
			}
			if path != testCase.ExpectedPath {
				t.Errorf("Expected %s, but acutually %s", testCase.ExpectedPath, path)
			}
		})
	}
}

func moduleFS() afero.Fs {
	fs := afero.NewMemMapFs()
	fs.MkdirAll("/atgen/lib", os.ModePerm)
	file, _ := fs.Create("/atgen/go.mod")
	defer file.Close()
	io.WriteString(file, testGoMod)
	return fs
}

func gopathFS() afero.Fs {
	fs := afero.NewMemMapFs()
	fs.MkdirAll("/go/src/github.com/aktsk/atgen", os.ModePerm)
	return fs
}
