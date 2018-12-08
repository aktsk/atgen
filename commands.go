package main

import (
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
			Name:  "yamlDir",
			Value: ".",
			Usage: "directory that has test definition yaml files",
		},
		cli.StringFlag{
			Name:  "templateFile",
			Value: "template_test.go",
			Usage: "template file defines test code",
		},
		cli.StringFlag{
			Name:  "dir",
			Value: ".",
			Usage: "output directory to write generated test files",
		},
	},
}

func doGen(c *cli.Context) error {
	files, err := filepath.Glob(filepath.Join(c.String("yamlDir"), "*.y*ml"))
	if err != nil {
		return err
	}

	for _, file := range files {
		generator := atgen.Generator{
			Yaml:     file,
			Template: c.String("templateFile"),
			Dir:      c.String("dir"),
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
