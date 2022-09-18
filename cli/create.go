package cli

import (
	"encoding/json"
	"fmt"
	"nemo/build"
	"os"
)

const StartMessage = `==----------==
New workspace is created.
Before you start, You must install skin to /skin directory.
Read guide about setting up skin for your workspace.`

func createNewSpace() {
	wd, _ := os.Getwd()
	_, postexist := os.Stat("post")
	_, skinexist := os.Stat("skin")
	_, maniexist := os.Stat("manifest.json")

	if !(os.IsNotExist(postexist) && os.IsNotExist(skinexist) && os.IsNotExist(maniexist)) {
		fmt.Println("workspace is already exist")
		return
	}

	fmt.Println("? Name:")
	bname := Prompt(true)

	fmt.Println("? Author:")
	bauthor := Prompt(false)

	fmt.Println("? Language:")
	blang := Prompt(false)

	fmt.Println("? Repository (optional, need for publish):")
	brepo := Prompt(false)

	os.Mkdir("post", os.ModePerm)
	os.Mkdir("skin", os.ModePerm)

	manifest := build.Manifest{
		Name:   bname,
		Author: bauthor,
		Lang:   blang,
		Repo:   brepo,
	}

	manibuild, err := json.Marshal(manifest)
	if err != nil {
		panic(err)
	}

	ferr := os.WriteFile("manifest.json", manibuild, 0777)
	if ferr != nil {
		panic(ferr)
	}

	os.WriteFile((wd + "/post/about.ps"), []byte("Write about message here."), 0777)
	os.Chdir("post")
	os.Mkdir("res", os.ModePerm)

	fmt.Println("Build Complete")
	fmt.Println(StartMessage)
}
