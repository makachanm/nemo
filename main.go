package main

import (
	"nemo/cli"
	"nemo/utils"
	"os"
)

var (
	BuildDate, Arch, Version string
)

func main() {
	vinfo := utils.MakeVersionInfo(BuildDate, Arch, Version)
	app := cli.MakeCli(vinfo)
	app.Handle(os.Args)
}
