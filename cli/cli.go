package cli

import (
	"fmt"
	"nemo/build"
)

const (
	BuildCommand   = "build"
	NewPostCommand = "newpost"
	CreateCommand  = "create"
)

type Interface struct {
	Usage      string
	commandMap map[string]func()
}

func MakeCli() Interface {
	ci := Interface{
		Usage: `Usage: nemo <command>

Commands:
  build       Build the site
  newpost     Create a new post
  create      Create a new space`,

		commandMap: map[string]func(){
			BuildCommand:   buildHandler,
			NewPostCommand: GeneratePost,
			CreateCommand:  createNewSpace,
		},
	}

	return ci
}

func (ci *Interface) Handle(args []string) {
	if len(args) < 2 {
		fmt.Println(ci.Usage)
		return
	}

	command := args[1]

	handler, ok := ci.commandMap[command]
	if !ok {
		fmt.Printf("\x1b[31merror\x1b[0m: unknown command '%s'\n\n", command)
		fmt.Println(ci.Usage)
		return
	}

	handler()
}

func buildHandler() {
	var b build.Builder

	build.MakeNewBuilder(&b)

	b.Build()
}
