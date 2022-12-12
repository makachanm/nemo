package main

import (
	"nemo/cli"
	"nemo/utils"
	"os"
)

var (
	BuildDate, Arch string
)

func main() {
	vinfo := utils.MakeVersionInfo(BuildDate, Arch)
	app := cli.MakeCli(vinfo)
	app.Handle(os.Args)
}
