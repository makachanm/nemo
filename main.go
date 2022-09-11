package main

import (
	"nemo/cli"
	"os"
)

func main() {
	app := cli.MakeCli()
	app.Handle(os.Args)
	//os.Args
}
