package main

import (
	"io"
	"os"
	"path/filepath"

	atgen "github.com/mizzy/atgen/lib"
	"github.com/urfave/cli"
)

var Commands = []cli.Command{
	commandGen,
	//commandDiff,
}

var commandGen = cli.Command{
	Name:  "gen",
	Usage: "Generate test code",
	Description: `
    Generete test code according to yaml and template.
`,
	Action: doGen,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "templateDir",
			Value: ".",
			Usage: "template directory that has template yaml and code",
		},
		cli.StringFlag{
			Name:  "outputDir",
			Value: ".",
			Usage: "output directory to write generated test files",
		},
	},
}

func doGen(c *cli.Context) error {
	templateDir := c.String("templateDir")
	outputDir := c.String("outputDir")

	testFiles, err := filepath.Glob(filepath.Join(templateDir, "*_test.go"))
	if err != nil {
		return err
	}

	for _, testFile := range testFiles {
		if testFile != "template_test.go" {
			src := filepath.Join(templateDir, testFile)
			dest := filepath.Join(outputDir, testFile)
			err := copyFile(src, dest)
			if err != nil {
				return err
			}
		}
	}

	yamlFiles, err := filepath.Glob(filepath.Join(templateDir, "*.y*ml"))
	if err != nil {
		return err
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
			return err
		}

		err = generator.Generate()
		if err != nil {
			return err
		}
	}

	return nil
}

func copyFile(s, d string) error {
	src, err := os.Open(s)
	if err != nil {
		return err
	}
	defer src.Close()

	dest, err := os.Create(d)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, src)
	if err != nil {
		return err
	}

	return nil
}
