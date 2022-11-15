package main

import (
	"fmt"
	"nemo/cli"
	"os"
)

var (
	Version, BuildDate string
)

func init() {
	fmt.Printf("NEMO Builder version %s - %s \n", Version, BuildDate)
}

func main() {
	app := cli.MakeCli()
	app.Handle(os.Args)
}
