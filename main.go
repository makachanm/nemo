package main

import (
	"fmt"
	"nemo/cli"
	"nemo/utils"
	"os"
)

var (
	BuildDate, Arch, Version string
)

func main() {
	vinfo := utils.MakeVersionInfo(BuildDate, Arch, Version)
	config, err := utils.LoadConfig()
	if err != nil {
		fmt.Println("Error in loading configuration: ", err)
		os.Exit(-1)
	}

	app := cli.MakeCli(vinfo, config)
	app.Handle(os.Args)
}
