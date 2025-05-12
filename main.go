package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var version = "0.9.0"

func main() {
	ctx := context.Background()
	err := newApp().Run(ctx, os.Args)
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func newApp() *cli.Command {
	var app cli.Command
	app.Name = "atgen"
	app.Usage = "Generate API test code from Request/Response definition or show diff between Request/Response definition and API definition"
	app.Version = version
	app.Authors = []any{
		"Akatsuki Inc.",
	}
	app.Commands = commands
	return &app
}
