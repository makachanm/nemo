package cli

import (
	"fmt"
	"nemo/build"
)

type CliInterface struct {
}

func MakeCli() CliInterface {
	return CliInterface{}
}

func (ci *CliInterface) Handle(args []string) {
	if len(args) < 2 {
		fmt.Println("not enough arguments")
		return
	}

	switch args[1] {
	case "build":
		build_handler()
	}
}

func build_handler() {
	b := build.MakeNewBuilder()
	b.Build()
}
