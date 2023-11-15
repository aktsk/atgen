package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var version = "0.8.0"

func main() {
	err := newApp().Run(os.Args)
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "atgen"
	app.Usage = "Generate API test code from Request/Response definition or show diff between Request/Response definition and API definition"
	app.Version = version
	app.Authors = []*cli.Author{
		{Name: "Akatsuki Inc."},
	}
	app.Commands = commands
	return app
}
