package main

import (
	"nemo/cli"
	"os"
)

/*
var (
	Version, BuildDate string
)
*/

func main() {
	app := cli.MakeCli()
	app.Handle(os.Args)
}
