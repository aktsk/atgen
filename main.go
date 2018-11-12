package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var Version string = "0.1.0"

func main() {
	err := newApp().Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "atgen"
	app.Usage = "Generate API test code from Request/Response definition or show diff between Request/Response definition and API definition"
	app.Version = Version
	app.Author = "Gosuke Miyashita"
	app.Email = "gosukenator@gmail.com"
	app.Commands = Commands
	return app
}
