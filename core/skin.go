package build

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
)

type SkinPath struct {
	Index  string `json:"index"`
	Post   string `json:"post"`
	About  string `json:"about"`
	Header string `json:"header"`
	Footer string `json:"footer"`
	Nav    string `json:"nav"`
	Tags   string `json:"tags"`
}

type SkinConfig struct {
	IndexNum int    `json:"indexnum"`
	DateType string `json:"datetype"`
}

type SkinInfo struct {
	Name    string     `json:"name"`
	Author  string     `json:"author"`
	Summary string     `json:"summary"`
	Repo    string     `json:"repo"`
	Conf    SkinConfig `json:"config"`
	Paths   SkinPath   `json:"paths"`
}

type Skin struct {
	Info SkinInfo
}

func MakeSkin() Skin {
	return Skin{}
}

func (s *Skin) GetSkin() error {
	skinpath, perr := os.Getwd()

	if perr != nil {
		return perr
	}

	//_, skinexist := os.Stat(skinpath + "/skin/skin.json")
	_, skinexist := os.Stat(filepath.Join(skinpath, "skin", "skin.json"))
	if os.IsNotExist(skinexist) {
		return errors.New("skin is not exist")
	}

	ctx, ferr := os.ReadFile(filepath.Join(skinpath, "skin", "skin.json"))

	if ferr != nil {
		return ferr
	}

	var sinfo = SkinInfo{}
	jerr := json.Unmarshal(ctx, &sinfo)

	if jerr != nil {
		return jerr
	}

	s.Info = sinfo

	var infos = SkinInfo{}

	ref := reflect.ValueOf(&sinfo.Paths).Elem()
	sref := reflect.ValueOf(&infos.Paths).Elem()

	for i := 0; i < ref.NumField(); i++ {
		pelmVal := ref.Field(i)
		pelmTyp := ref.Type().Field(i)

		srefval := sref.FieldByName(pelmTyp.Name)
		pval := skinpath + (pelmVal.Interface().(string))
		srefval.SetString(pval)
	}

	s.Info.Paths = infos.Paths

	return nil
}
