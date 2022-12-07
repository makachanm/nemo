package cli

import (
	"fmt"
	"nemo/build"
)

type Interface struct {
	Usage string
}

func MakeCli() Interface {
	return Interface{
		Usage: "Usage: nemo <command>\n\n" +
			"Commands:\n" +
			"  build       Build the site\n" +
			"  newpost     Create a new post\n" +
			"  create      Create a new space\n",
	}
}

func (ci *Interface) Handle(args []string) {
	if len(args) < 2 {
		fmt.Println(ci.Usage)
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
		fmt.Println("\x1b[31merror\x1b[0m: unknown command")
		fmt.Println("\n", ci.Usage)
	}
}

func buildHandler() {
	b := build.MakeNewBuilder()
	b.Build()
}
