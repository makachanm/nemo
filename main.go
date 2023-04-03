package main

import (
	"log"
	"nemo/cli"
	"nemo/utils"
	"os"
)

var (
	Version   string
	Arch      string
	BuildDate string
)

func main() {
	vinfo := utils.MakeVersionInfo(BuildDate, Arch, Version)
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatal("Error in loading configuration: ", err)
	}

	app := cli.MakeCli(vinfo, config)
	app.Handle(os.Args)
}
