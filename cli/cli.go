package cli

import (
	"fmt"
	build "nemo/core"
	"nemo/utils"
)

const UsageT = `Usage: nemo <command>

Commands:
  build       Build the site
  newpost     Create a new post
  create      Create a new space
  help        Print this help message
`

const (
	BuildCommand    = "build"
	NewPostCommand  = "newpost"
	CreateCommand   = "create"
	ShowHelpMessage = "help"
)

type Interface struct {
	Usage       string
	commandMap  map[string]func()
	versionInfo utils.VersionInfo
	configInfo  utils.Config
}

func MakeCli(vinfo utils.VersionInfo, config utils.Config) Interface {
	ci := Interface{
		Usage: UsageT,

		versionInfo: vinfo,
		configInfo:  config,
	}

	ci.commandMap = map[string]func(){
		BuildCommand:    ci.buildHandler,
		NewPostCommand:  GeneratePost,
		CreateCommand:   createNewSpace,
		ShowHelpMessage: ci.printHelpMessage,
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

func (i Interface) buildHandler() {
	var b build.Builder

	build.MakeNewBuilder(&b, i.versionInfo, i.configInfo)

	b.Build()
}

func (i Interface) printHelpMessage() {
	fmt.Print(UsageT)
}
