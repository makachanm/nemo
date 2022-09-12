package cli

import (
	"fmt"
	"nemo/build"
)

type Interface struct {
}

func MakeCli() Interface {
	return Interface{}
}

func (ci *Interface) Handle(args []string) {
	if len(args) < 2 {
		fmt.Println("not enough arguments")
		return
	}

	switch args[1] {
	case "build":
		buildHandler()

	case "newpost":
		GeneratePost()

	case "create":
		createNewSpace()

	default:
		fmt.Println("unknown command")
	}
}

func buildHandler() {
	b := build.MakeNewBuilder()
	b.Build()
}
