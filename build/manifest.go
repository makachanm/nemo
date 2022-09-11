package build

import (
	"encoding/json"
	"fmt"
	"os"
)

type Manifest struct {
	Name   string `json:"name"`
	Lang   string `json:"lang"`
	Author string `json:"author"`
	Repo   string `json:"repository"`
}

func GetManifest() Manifest {
	wd, perr := os.Getwd()

	if perr != nil {
		panic(perr)
	}

	_, maniexist := os.Stat(wd + "/manifest.json")
	if os.IsNotExist(maniexist) {
		fmt.Println("Manifest is not exist")
		os.Exit(1)
	}

	ctx, ferr := os.ReadFile(wd + "/manifest.json")

	if ferr != nil {
		panic(ferr)
	}

	var sinfo = Manifest{}
	jerr := json.Unmarshal(ctx, &sinfo)

	if jerr != nil {
		panic(jerr)
	}

	return sinfo
}
