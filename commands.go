package main

import (
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
			Name:  "yaml",
			Value: "tests.yml",
			Usage: "yaml file defines requests and responses for testing",
		},
		cli.StringFlag{
			Name:  "template",
			Value: "template.go",
			Usage: "template file defines test code",
		},
	},
}

func doGen(c *cli.Context) error {
	generator := atgen.Generator{
		Yaml:     c.String("yaml"),
		Template: c.String("template"),
	}

	err := generator.ParseYaml()
	if err != nil {
		return err
	}

	err = generator.Generate()
	if err != nil {
		return err
	}

	return err
}
