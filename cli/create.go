package cli

import (
	"encoding/json"
	"fmt"
	build "nemo/core"
	"nemo/utils"
	"os"
)

const StartMessage = `==----------==
New workspace is created.
Before you start, You must install skin to /skin directory.
Read guide about setting up skin for your workspace.`

func createNewSpace(vinfo utils.VersionInfo) {
	wd, _ := os.Getwd()
	_, maniexist := os.Stat("manifest.json")

	if !os.IsNotExist(maniexist) {
		fmt.Println("workspace is already exist")
		return
	}

	fmt.Println("? Name:")
	bname, err := Prompt(true)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("? Author:")
	bauthor, err := Prompt(false)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("? Language:")
	blang, err := Prompt(false)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("? Repository (optional, need for publish):")
	brepo, err := Prompt(false)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("? Domain (optional, need for generate RSS Feed):")
	bdom, err := Prompt(false)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.Mkdir("post", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating post directory:", err)
		return
	}

	err = os.Mkdir("skin", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating skin directory:", err)
		return
	}

	manifest := build.Manifest{
		Name:   bname,
		Author: bauthor,
		Lang:   blang,
		Repo:   brepo,
		Domain: bdom,
	}

	manibuild, err := json.Marshal(manifest)
	if err != nil {
		fmt.Println("Error marshalling manifest:", err)
		return
	}

	ferr := os.WriteFile("manifest.json", manibuild, 0777)
	if ferr != nil {
		fmt.Println("Error writing manifest file:", ferr)
		return
	}

	ferr = os.WriteFile(wd+"/post/about.ps", []byte("Write about message here."), 0777)
	if ferr != nil {
		fmt.Println("Error writing about.ps file:", ferr)
		return
	}

	err = os.Chdir("post")
	if err != nil {
		fmt.Println("Error changing directory:", err)
		return
	}

	err = os.Mkdir("res", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating res directory:", err)
		return
	}

	fmt.Println("Build Complete")
	fmt.Println(StartMessage)
}
