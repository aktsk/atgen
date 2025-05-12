package main

import (
	"context"
	"io"
	"os"
	"path/filepath"

	atgen "github.com/aktsk/atgen/lib"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"
)

var commands = []*cli.Command{
	commandGen,
	//commandDiff,
}

var commandGen = &cli.Command{
	Name:  "gen",
	Usage: "Generate test code",
	Description: `
    Generete test code according to yaml and template.
`,
	Action: doGen,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "templateDir",
			Aliases: []string{"t"},
			Value:   ".",
			Usage:   "template directory that has template yaml and code",
		},
		&cli.StringFlag{
			Name:    "outputDir",
			Aliases: []string{"o"},
			Value:   ".",
			Usage:   "output directory to write generated test files",
		},
	},
}

func doGen(_ context.Context, cmd *cli.Command) error {
	templateDir := cmd.String("templateDir")
	outputDir := cmd.String("outputDir")

	testFiles, err := filepath.Glob(filepath.Join(templateDir, "*_test.go"))
	if err != nil {
		return errors.WithStack(err)
	}

	for _, testFile := range testFiles {
		base := filepath.Base(testFile)
		if base != "template_test.go" {
			src := filepath.Join(templateDir, base)
			dest := filepath.Join(outputDir, base)
			err := copyFile(src, dest)
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}

	yamlFiles, err := filepath.Glob(filepath.Join(templateDir, "*.y*ml"))
	if err != nil {
		return errors.WithStack(err)
	}

	for _, yamlFile := range yamlFiles {
		generator := atgen.Generator{
			Yaml:        yamlFile,
			Template:    filepath.Join(templateDir, "template_test.go"),
			TemplateDir: templateDir,
			OutputDir:   outputDir,
		}

		err := generator.ParseYaml()
		if err != nil {
			return errors.WithStack(err)
		}

		err = generator.Generate()
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func copyFile(s, d string) error {
	src, err := os.Open(s)
	if err != nil {
		return errors.WithStack(err)
	}
	defer src.Close()
	dstdir := filepath.Dir(d)

	tmpdst, err := os.CreateTemp(dstdir, "tmp-")
	if err != nil {
		return errors.WithStack(err)
	}
	defer (func() {
		if tmpdst != nil {
			f := tmpdst.Name()
			_ = tmpdst.Close()
			_ = os.Remove(f)
		}
	})()

	_, err = io.Copy(tmpdst, src)
	if err != nil {
		return errors.WithStack(err)
	}
	if err = tmpdst.Close(); err != nil {
		return errors.WithStack(err)
	}
	if err = os.Rename(tmpdst.Name(), d); err != nil {
		return errors.WithStack(err)
	}
	tmpdst = nil
	return nil
}
