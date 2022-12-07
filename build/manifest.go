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

func GetManifest() (Manifest, error) {
	wd, err := os.Getwd()
	if err != nil {
		return Manifest{}, err
	}

	if _, err := os.Stat(wd + "/manifest.json"); os.IsNotExist(err) {
		return Manifest{}, fmt.Errorf("manifest file does not exist")
	}

	ctx, err := os.ReadFile(wd + "/manifest.json")
	if err != nil {
		return Manifest{}, err
	}

	var sinfo = Manifest{}
	if err := json.Unmarshal(ctx, &sinfo); err != nil {
		return Manifest{}, err
	}

	return sinfo, nil
}
