package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	UseLegacyParser bool `json:"useLegacyParser"`
}

func LoadConfig() (Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error in getting fs: ", err)
		return Config{}, err
	}

	if _, err := os.Stat(filepath.Join(wd, "config.json")); os.IsNotExist(err) {
		fmt.Println("configuration file does not exist. Creating new configuration file.")
		MakeConfig()
		return Config{}, errors.New("config file not exist")

	}

	ctx, err := os.ReadFile(filepath.Join(wd, "config.json"))
	if err != nil {
		return Config{}, err
	}

	var cfg = Config{}
	err = json.Unmarshal(ctx, &cfg)
	if err != nil {
		return Config{}, err
	}

	return Config{}, nil
}

func MakeConfig() {
	config := Config{
		UseLegacyParser: false,
	}

	d, err := json.Marshal(config)
	if err != nil {
		fmt.Println("Error in making config file: ", err)
		return
	}

	ferr := os.WriteFile("config.json", d, 0777)
	if ferr != nil {
		fmt.Println("Error in writing config file:", ferr)
		return
	}
}
