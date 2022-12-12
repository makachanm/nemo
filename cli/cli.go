package cli

import (
	"fmt"
	"nemo/build"
	"nemo/utils"
)

const (
	BuildCommand    = "build"
	NewPostCommand  = "newpost"
	CreateCommand   = "create"
	ShowHelpMessage = "help"
)

type Interface struct {
	Usage       string
	commandMap  map[string]func(utils.VersionInfo)
	versionInfo utils.VersionInfo
}

func MakeCli(vinfo utils.VersionInfo) Interface {
	ci := Interface{
		Usage: `Usage: nemo <command>

Commands:
  build       Build the site
  newpost     Create a new post
  create      Create a new space
  help        Print this help message
`,

		commandMap: map[string]func(utils.VersionInfo){
			BuildCommand:    buildHandler,
			NewPostCommand:  GeneratePost,
			CreateCommand:   createNewSpace,
			ShowHelpMessage: printHelpMessage,
		},

		versionInfo: vinfo,
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

	handler(ci.versionInfo)
}

func buildHandler(vinfo utils.VersionInfo) {
	var b build.Builder

	build.MakeNewBuilder(&b, vinfo)

	b.Build()
}

func printHelpMessage(vinfo utils.VersionInfo) {
	fmt.Println(MakeCli(vinfo).Usage)
}
