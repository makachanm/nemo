package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"nemo/build"
	"os"
)

func createNewSpace() {
	_, postexist := os.Stat("post")
	_, skinexist := os.Stat("skin")
	_, maniexist := os.Stat("manifest.json")

	if !(os.IsNotExist(postexist) && os.IsNotExist(skinexist) && os.IsNotExist(maniexist)) {
		fmt.Println("workspace is already exist")
		return
	}

	red := bufio.NewReader(os.Stdin)

	fmt.Println("? Name:")
	bname, _ := red.ReadString('\n')

	fmt.Println("? Author:")
	bauthor, _ := red.ReadString('\n')

	fmt.Println("? Language:")
	blang, _ := red.ReadString('\n')

	fmt.Println("? Repository (optional, need for publish):")
	brepo, _ := red.ReadString('\n')

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

	fmt.Println("Build Complete")
}
